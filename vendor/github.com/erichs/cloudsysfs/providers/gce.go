package providers

import (
	"io/ioutil"
	"strings"
)

func GCE(sysfscheck chan<- string) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/product_name")
	if err != nil {
		sysfscheck <- ""
	}
	if strings.Contains(string(data), "Google") {
		sysfscheck <- "gce"
	}
	sysfscheck <- ""
}
