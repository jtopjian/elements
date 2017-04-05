package providers

import (
	"io/ioutil"
	"strings"
)

func OpenStack(sysfscheck chan<- string) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if err != nil {
		sysfscheck <- ""
	}
	if strings.Contains(string(data), "OpenStack Foundation") {
		sysfscheck <- "openstack"
	}
	sysfscheck <- ""
}
