package providers

import (
	"io/ioutil"
	"strings"
)

func Azure(sysfscheck chan<- string) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
	if err != nil {
		sysfscheck <- ""
	}
	if strings.Contains(string(data), "Microsoft Corporation") {
		sysfscheck <- "azure"
	}
	sysfscheck <- ""
}
