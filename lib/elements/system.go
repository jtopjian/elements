package elements

import (
	"net"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net"
)

// GetSystemElements is the main function to call to collect all elements on the system.
// You shouldn't need to interact with anything other than this function.
func (e *Elements) GetSystemElements() (interface{}, error) {
	systemElements := make(map[string]interface{})

	v, err := getProcessors()
	if err != nil {
		return nil, err
	}
	systemElements["processors"] = v

	v, err = getDisk()
	if err != nil {
		return nil, err
	}
	systemElements["disks"] = v

	v, err = getHost()
	if err != nil {
		return nil, err
	}
	systemElements["host"] = v

	v, err = getMemory()
	if err != nil {
		return nil, err
	}
	systemElements["memory"] = v

	v, err = getInterfaces()
	if err != nil {
		return nil, err
	}
	systemElements["interfaces"] = v

	return systemElements, nil
}

func getProcessors() (interface{}, error) {
	v := make(map[string]interface{})

	// total processor count
	count, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	v["count"] = count

	// information about each processor
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	v["info"] = cpuInfo

	return v, nil
}

func getDisk() (interface{}, error) {
	return disk.Partitions(false)
}

func getHost() (interface{}, error) {
	/*
		hostInfo, err := host.Info()
		if err != nil {
			return nil, err
		}
	*/

	return host.Info()
}

func getMemory() (interface{}, error) {
	v := make(map[string]interface{})

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	v["virtual"] = vm

	sm, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	v["swap"] = sm

	return v, nil
}

func getInterfaces() (interface{}, error) {
	v := make(map[string]interface{})

	interfaces, err := psnet.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		iface := make(map[string]interface{})
		ipv4 := []map[string]interface{}{}
		ipv6 := []map[string]interface{}{}
		interfaceName := i.Name
		iface["mtu"] = i.MTU
		iface["hardwareaddr"] = i.HardwareAddr
		iface["flags"] = i.Flags

		for _, addr := range i.Addrs {
			ip, ipnet, err := net.ParseCIDR(addr.Addr)
			if err == nil {
				ipinfo := make(map[string]interface{})
				ipinfo["address"] = ip
				ipinfo["cidr"] = addr.Addr
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					ipv4 = append(ipv4, ipinfo)
				} else {
					ipv6 = append(ipv6, ipinfo)
				}
			}
		}

		if len(ipv4) > 0 {
			iface["ipv4"] = ipv4
		}

		if len(ipv6) > 0 {
			iface["ipv6"] = ipv6
		}

		v[interfaceName] = iface
	}

	return v, nil
}
