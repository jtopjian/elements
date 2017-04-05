# cloudsysfs
Detect Cloud Provider from DMI data on the local sysfs (/sys) filesystem

## Usage

```go
package main

import (
  "fmt"

  "github.com/erichs/cloudsysfs"
)

func main() {
  switch cloudsysfs.Detect() {
  case "aws":
    fmt.Println("Amazon Web Services")
  case "azure":
    fmt.Println("Microsofte Azure")
  case "digitalocean":
    fmt.Println("Digital Ocean")
  case "gce":
    fmt.Println("Google Compute Engine")
  case "openstack":
    fmt.Println("OpenStack")
  default:
    fmt.Println("No cloud detected")
  }
}

``` 

## Motivation

Inspired by [cloudid](https://github.com/appscode/cloudid), I wanted a simple and fast mechanism for determining if the local environment is a cloud I want to support. [cloudid](https://github.com/appscode/cloudid) is great, but takes the (admittedly more robust) approach of fingerprinting metadata signatures from APIPA HTTP API endpoints (`http://169.254.x.x/latest/metadata` and friends). 

This library takes the approach of attempting to read from the local [sysfs filesystem](https://en.wikipedia.org/wiki/Sysfs), looking for unambiguous vendor or product files that identify the cloud provider. The /sys filesystem is provided by the Linux kernel, so only Linux is currently supported.

## Supported Operating Systems

Linux

## Supported Cloud Providers
| provider_id | Name                  
|-------------|-----------
|aws          | Amazon Web Services  
|azure        | Microsoft Azure      
|digitalocean | DigitalOcean          
|gce          | Google Compute Engine
|openstack    | OpenStack

## Example

```      
package main

import (
  "fmt"

  "github.com/erichs/cloudsysfs"
)

func main() {
  switch cloudsysfs.Detect() {
  case "aws":
    fmt.Println("Amazon Web Services")
  case "azure":
    fmt.Println("Microsofte Azure")
  case "digitalocean":
    fmt.Println("Digital Ocean")
  case "gce":
    fmt.Println("Google Compute Engine")
  case "openstack":
    fmt.Println("OpenStack")
  default:
    fmt.Println("No cloud detected")
  }
}
```
