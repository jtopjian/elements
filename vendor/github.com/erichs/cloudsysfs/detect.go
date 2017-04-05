package cloudsysfs

import (
	"github.com/erichs/cloudsysfs/providers"
)

var Providers = [...]func(chan<- string){
	providers.AWS,
	providers.Azure,
	providers.DigitalOcean,
	providers.GCE,
	providers.OpenStack,
}

func Detect() string {
	sysfscheck := make(chan string)
	for _, cloud := range Providers {
		go cloud(sysfscheck)
	}

	provider := ""
	for i := 0; i < len(Providers); i++ {
		v := <-sysfscheck
		if v != "" {
			provider = v
		}
	}

	return provider
}
