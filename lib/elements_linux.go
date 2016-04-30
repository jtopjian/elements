// +build linux

package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jtopjian/elements/utils"
	"golang.org/x/sys/unix"
)

// Constants
const (
	// LINUX_SYSINFO_LOADS_SCALE has been described elsewhere as a "magic" number.
	// It reverts the calculation of "load << (SI_LOAD_SHIFT - FSHIFT)" done in the original load calculation.
	LINUX_SYSINFO_LOADS_SCALE float64 = 65536.0
)

// GetSystemElements is the main function to call to collect all elements on the system.
// You shouldn't need to interact with anything other than this function.
func (e *Elements) GetSystemElements() *SystemElements {
	elements := new(SystemElements)
	var wg sync.WaitGroup

	wg.Add(11)
	go elements.getOSRelease(&wg)
	go elements.getInterfaces(&wg)
	go elements.getBootID(&wg)
	go elements.getMachineID(&wg)
	go elements.getUname(&wg)
	go elements.getSysInfo(&wg)
	go elements.getDate(&wg)
	go elements.getFileSystems(&wg)
	go elements.getDMI(&wg)
	go elements.getBlockDevices(&wg)
	go elements.getProcessors(&wg)

	wg.Wait()
	return elements
}

// SystemElements holds the system elements.
type SystemElements struct {
	Architecture string
	BootID       string
	DMI          DMI
	Date         Date
	Domainname   string
	Hostname     string
	Network      Network
	Kernel       Kernel
	MachineID    string
	Memory       Memory
	OSRelease    OSRelease
	Swap         Swap
	Uptime       int64
	LoadAverage  LoadAverage
	FileSystems  FileSystems
	BlockDevices BlockDevices
	Processors   Processors

	mu sync.Mutex
}

// DMI holds the DMI / Hardware Information.
type DMI struct {
	BIOSDate        string
	BIOSVendor      string
	BIOSVersion     string
	ChassisAssetTag string
	ChassisSerial   string
	ChassisType     string
	ChassisVendor   string
	ChassisVersion  string
	ProductName     string
	ProductSerial   string
	ProductUUID     string
	ProductVersion  string
	SysVendor       string
}

// Holds the load average elements.
type LoadAverage struct {
	One  string
	Five string
	Ten  string
}

// Date holds the date elements.
type Date struct {
	Unix     int64
	UTC      string
	Timezone string
	Offset   int
}

// Swap holds the swap elements.
type Swap struct {
	Total uint64
	Free  uint64
}

// OSRelease holds the OS release elements.
type OSRelease struct {
	Name       string
	ID         string
	PrettyName string
	Version    string
	VersionID  string
	CodeName   string
}

// Kernel holds the kernel elements.
type Kernel struct {
	Name    string
	Release string
	Version string
}

// Memory holds the memory elements.
type Memory struct {
	Total    uint64
	Free     uint64
	Shared   uint64
	Buffered uint64
}

// Network holds the network elements.
type Network struct {
	Interfaces Interfaces
}

// Interfaces holds the interface (NIC) elements.
type Interfaces map[string]Interface

// Interface holds elements for a single interface (NIC).
type Interface struct {
	Name         string
	Index        int
	HardwareAddr string
	IpAddresses  []string
	Ip4Addresses []Ip4Address
	Ip6Addresses []Ip6Address
}

type Ip4Address struct {
	CIDR    string
	Ip      string
	Netmask string
}

type Ip6Address struct {
	CIDR   string
	Ip     string
	Prefix int
}

// FileSystems holds the Filesystem elements.
type FileSystems map[string]FileSystem

// FileSystem holds elements for a filesystem (man fstab).
type FileSystem struct {
	Device     string
	MountPoint string
	Type       string
	Options    []string
	DumpFreq   uint64
	PassNo     uint64
}

// BlockDevices holds the BlockDevice elements.
type BlockDevices map[string]BlockDevice

// BlockDevice holds elements for a block device
type BlockDevice struct {
	Device       string
	Size         uint64
	Vendor       string
	ReadIOs      uint64
	ReadMerges   uint64
	ReadSectors  uint64
	ReadTicks    uint64
	WriteIOs     uint64
	WriteMerges  uint64
	WriteSectors uint64
	WriteTicks   uint64
	InFlight     uint64
	IOTicks      uint64
	TimeInQueue  uint64
}

// Processors holds elements about the Processors / CPUs.
type Processors struct {
	Count     int
	Processor []Processor
}

// Processor holds elements about a single Processor / CPU.
type Processor struct {
	VendorID  string
	CPUFamily uint64
	Model     uint64
	ModelName string
	MHz       string
	CacheSize string
	CPUCores  uint64
	Flags     []string
	BogoMips  float64
}

func (f *SystemElements) getDate(wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now()
	f.Date.Unix = now.Unix()
	f.Date.UTC = now.UTC().String()
	f.Date.Timezone, f.Date.Offset = now.Zone()

	return
}

func (f *SystemElements) getSysInfo(wg *sync.WaitGroup) error {
	defer wg.Done()

	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.Memory.Total = info.Totalram
	f.Memory.Free = info.Freeram
	f.Memory.Shared = info.Sharedram
	f.Memory.Buffered = info.Bufferram

	f.Swap.Total = info.Totalswap
	f.Swap.Free = info.Freeswap

	f.Uptime = info.Uptime

	f.LoadAverage.One = fmt.Sprintf("%.2f", float64(info.Loads[0])/LINUX_SYSINFO_LOADS_SCALE)
	f.LoadAverage.Five = fmt.Sprintf("%.2f", float64(info.Loads[1])/LINUX_SYSINFO_LOADS_SCALE)
	f.LoadAverage.Ten = fmt.Sprintf("%.2f", float64(info.Loads[2])/LINUX_SYSINFO_LOADS_SCALE)

	return nil
}

func (f *SystemElements) getOSRelease(wg *sync.WaitGroup) error {
	defer wg.Done()
	osReleaseFile, err := os.Open("/etc/os-release")
	if err != nil {
		return err
	}
	defer osReleaseFile.Close()

	f.mu.Lock()
	defer f.mu.Unlock()
	scanner := bufio.NewScanner(osReleaseFile)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), "=")
		if len(columns) > 1 {
			key := columns[0]
			value := strings.Trim(strings.TrimSpace(columns[1]), `"`)
			switch key {
			case "NAME":
				f.OSRelease.Name = value
			case "ID":
				f.OSRelease.ID = value
			case "PRETTY_NAME":
				f.OSRelease.PrettyName = value
			case "VERSION":
				f.OSRelease.Version = value
			case "VERSION_ID":
				f.OSRelease.VersionID = value
			}
		}
	}

	lsbFile, err := os.Open("/etc/lsb-release")
	if err != nil {
		return err
	}
	defer lsbFile.Close()

	scanner = bufio.NewScanner(lsbFile)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), "=")
		if len(columns) > 1 {
			key := columns[0]
			value := strings.Trim(strings.TrimSpace(columns[1]), `"`)
			switch key {
			case "DISTRIB_CODENAME":
				f.OSRelease.CodeName = value
			}
		}
	}

	return nil
}

func (f *SystemElements) getMachineID(wg *sync.WaitGroup) error {
	defer wg.Done()
	machineIDFile, err := os.Open("/etc/machine-id")
	if err != nil {
		return err
	}
	defer machineIDFile.Close()
	data, err := ioutil.ReadAll(machineIDFile)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.MachineID = strings.TrimSpace(string(data))
	return nil
}

func (f *SystemElements) getBootID(wg *sync.WaitGroup) error {
	defer wg.Done()
	bootIDFile, err := os.Open("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return err
	}
	defer bootIDFile.Close()
	data, err := ioutil.ReadAll(bootIDFile)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.BootID = strings.TrimSpace(string(data))
	return nil
}

func (f *SystemElements) getInterfaces(wg *sync.WaitGroup) error {
	defer wg.Done()
	ls, err := net.Interfaces()
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	m := make(Interfaces)
	for _, i := range ls {
		ipaddreses := make([]string, 0)
		ip4addrs := make([]Ip4Address, 0)
		ip6addrs := make([]Ip6Address, 0)

		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, ip := range addrs {
			cidr := ip.String()
			ipaddreses = append(ipaddreses, cidr)
			ip, ipnet, _ := net.ParseCIDR(cidr)
			if ip.To4() != nil {
				ip4addrs = append(ip4addrs, Ip4Address{cidr, ip.String(), toNetmask(ipnet.Mask)})
				continue
			}
			if ip.To16() != nil {
				ones, _ := ipnet.Mask.Size()
				ip6addrs = append(ip6addrs, Ip6Address{cidr, ip.String(), ones})
			}
		}
		m[i.Name] = Interface{
			Name:         i.Name,
			Index:        i.Index,
			HardwareAddr: i.HardwareAddr.String(),
			IpAddresses:  ipaddreses,
			Ip4Addresses: ip4addrs,
			Ip6Addresses: ip6addrs,
		}
	}
	f.Network.Interfaces = m

	return nil
}

func toNetmask(m net.IPMask) string {
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func (f *SystemElements) getUname(wg *sync.WaitGroup) error {
	defer wg.Done()

	var buf unix.Utsname
	err := unix.Uname(&buf)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.Domainname = utils.CharsToString(buf.Domainname)
	f.Architecture = utils.CharsToString(buf.Machine)
	f.Hostname = utils.CharsToString(buf.Nodename)
	f.Kernel.Name = utils.CharsToString(buf.Sysname)
	f.Kernel.Release = utils.CharsToString(buf.Release)
	f.Kernel.Version = utils.CharsToString(buf.Version)

	return nil
}

func (f *SystemElements) getFileSystems(wg *sync.WaitGroup) error {
	defer wg.Done()

	mtab, err := ioutil.ReadFile("/etc/mtab")
	if err != nil {
		return err
	}

	fsMap := make(FileSystems)

	f.mu.Lock()
	defer f.mu.Unlock()

	s := bufio.NewScanner(bytes.NewBuffer(mtab))
	for s.Scan() {
		line := s.Text()
		if string(line[0]) == "#" {
			continue
		}
		fields := strings.Fields(s.Text())
		fs := FileSystem{}
		fs.Device = fields[0]
		fs.MountPoint = fields[1]
		fs.Type = fields[2]
		fs.Options = strings.Split(fields[3], ",")
		dumpFreq, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return err
		}
		fs.DumpFreq = dumpFreq

		passNo, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return err
		}
		fs.PassNo = passNo

		fsMap[fs.Device] = fs
	}

	f.FileSystems = fsMap

	return nil
}

func (f *SystemElements) getDMI(wg *sync.WaitGroup) error {
	defer wg.Done()
	f.mu.Lock()
	defer f.mu.Unlock()

	var err error
	if f.DMI.BIOSDate, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/bios_date"); err != nil {
		return err
	}

	if f.DMI.BIOSVendor, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/bios_vendor"); err != nil {
		return err
	}

	if f.DMI.BIOSVersion, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/bios_version"); err != nil {
		return err
	}

	if f.DMI.ChassisAssetTag, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/chassis_asset_tag"); err != nil {
		return err
	}

	if f.DMI.ChassisSerial, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/chassis_serial"); err != nil {
		return err
	}

	if f.DMI.ChassisVendor, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/chassis_vendor"); err != nil {
		return err
	}

	if f.DMI.ChassisVersion, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/chassis_version"); err != nil {
		return err
	}

	if f.DMI.ProductName, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/product_name"); err != nil {
		return err
	}

	if f.DMI.ProductSerial, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/product_serial"); err != nil {
		return err
	}

	if f.DMI.ProductUUID, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/product_uuid"); err != nil {
		return err
	}

	if f.DMI.ProductVersion, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/product_version"); err != nil {
		return err
	}

	if f.DMI.SysVendor, err = utils.ReadFileAndReturnValue("/sys/class/dmi/id/sys_vendor"); err != nil {
		return err
	}

	return nil
}

func (f *SystemElements) getBlockDevices(wg *sync.WaitGroup) error {
	defer wg.Done()

	d, err := os.Open("/sys/block")
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		return err
	}

	bdMap := make(BlockDevices)
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, fi := range files {
		path := fmt.Sprintf("/sys/block/%s/device", fi.Name())
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			bd := BlockDevice{}
			bd.Device = fi.Name()

			sizePath := fmt.Sprintf("/sys/block/%s/size", fi.Name())
			size, err := utils.ReadFileAndReturnValue(sizePath)
			if err != nil {
				return err
			}

			bd.Size, _ = strconv.ParseUint(size, 10, 64)

			vendorPath := fmt.Sprintf("/sys/block/%s/device/vendor", fi.Name())
			if bd.Vendor, err = utils.ReadFileAndReturnValue(vendorPath); err != nil {
				return err
			}

			statPath := fmt.Sprintf("/sys/block/%s/stat", fi.Name())
			sf, err := os.Open(statPath)
			if err != nil {
				return err
			}
			defer sf.Close()

			scanner := bufio.NewScanner(sf)
			for scanner.Scan() {
				columns := strings.Fields(scanner.Text())
				if len(columns) == 11 {
					bd.ReadIOs, _ = strconv.ParseUint(columns[0], 10, 64)
					bd.ReadMerges, _ = strconv.ParseUint(columns[1], 10, 64)
					bd.ReadSectors, _ = strconv.ParseUint(columns[2], 10, 64)
					bd.ReadTicks, _ = strconv.ParseUint(columns[3], 10, 64)
					bd.WriteIOs, _ = strconv.ParseUint(columns[4], 10, 64)
					bd.WriteMerges, _ = strconv.ParseUint(columns[5], 10, 64)
					bd.WriteSectors, _ = strconv.ParseUint(columns[6], 10, 64)
					bd.WriteTicks, _ = strconv.ParseUint(columns[7], 10, 64)
					bd.InFlight, _ = strconv.ParseUint(columns[8], 10, 64)
					bd.IOTicks, _ = strconv.ParseUint(columns[9], 10, 64)
					bd.TimeInQueue, _ = strconv.ParseUint(columns[10], 10, 64)
				}
			}
			bdMap[bd.Device] = bd
		}
	}

	f.BlockDevices = bdMap

	return nil
}

func (f *SystemElements) getProcessors(wg *sync.WaitGroup) error {
	defer wg.Done()
	processorFile, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return err
	}
	defer processorFile.Close()

	f.mu.Lock()
	defer f.mu.Unlock()

	var cpuCount int = 0
	procs := []Processor{}
	p := Processor{}

	scanner := bufio.NewScanner(processorFile)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), ":")
		if len(columns) > 1 {
			key := strings.TrimSpace(columns[0])
			value := strings.TrimSpace(columns[1])

			switch key {
			case "processor":
				cpuCount += 1
			case "vendor_id":
				p.VendorID = value
			case "cpu family":
				p.CPUFamily, _ = strconv.ParseUint(value, 10, 64)
			case "model":
				p.Model, _ = strconv.ParseUint(value, 10, 64)
			case "model name":
				p.ModelName = value
			case "cpu MHz":
				p.MHz = value
			case "cache size":
				p.CacheSize = value
			case "cpu cores":
				p.CPUCores, _ = strconv.ParseUint(value, 10, 64)
			case "flags":
				value := strings.Fields(columns[1])
				p.Flags = value
			case "bogomips":
				p.BogoMips, _ = strconv.ParseFloat(value, 64)
			}
		} else {
			procs = append(procs, p)
			p = Processor{}
		}
	}

	f.Processors.Count = cpuCount
	f.Processors.Processor = procs

	return nil
}
