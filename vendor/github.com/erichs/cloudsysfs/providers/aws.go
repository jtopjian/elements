package providers

import (
	"io/ioutil"
	"strings"
)

func AWS(sysfscheck chan<- string) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/product_version")
	if err != nil {
		sysfscheck <- ""
	}
	if strings.Contains(string(data), "amazon") {
		sysfscheck <- "aws"
	}
	sysfscheck <- ""
}
