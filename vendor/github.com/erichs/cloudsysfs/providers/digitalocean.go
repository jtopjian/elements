package providers

import (
	"io/ioutil"
	"strings"
)

func DigitalOcean(sysfscheck chan<- string) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if err != nil {
		sysfscheck <- ""
	}
	if strings.Contains(string(data), "DigitalOcean") {
		sysfscheck <- "digitalocean"
	}
	sysfscheck <- ""
}
