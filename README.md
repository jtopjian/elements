# Elements

Elements retrieves information about a system such as CPU/Processors, disks, memory, and network interfaces. It is meant to be used to poll a system about it's static attributes.

This tool is similar to Facter, Ohai, Ansible facts, etc.

For dynamic attributes, such as load average and IOPS, use a more suitable metric collection tool.

## Install

Download the latest binary [release](https://github.com/jtopjian/elements/releases). Linux, Mac, FreeBSD, and Windows binaries are currently available, though testing has only been done on Linux and Mac.

## Usage

Run `elements` with no arguments on the command line to see the usage and options.

Elements is able to detect a standard set of information about a system. Simply run `elements get` on the command line and see the result:

```shell
$ elements get

{
  "system": {
    "host": {
			"bootTime": 1480266892,
			"hostid": "2DAAF9F7-0ED2-4DA7-BDD7-7A03F37C091A",
			"hostname": "jtdev",
			"kernelVersion": "4.4.0-45-generic",
			"os": "linux",
			"platform": "ubuntu",
			"platformFamily": "debian",
			"platformVersion": "16.04",
			"procs": 682,
			"uptime": 1678017,
			"virtualizationRole": "host",
			"virtualizationSystem": "kvm"
    },
		...
```

To retrieve a subset of elements, "walk" the tree using a dotted notation:

```shell
$ elements get -p system.interfaces.ens3

{
  "flags": [
   "up",
   "broadcast",
   "multicast"
  ],
  "hardwareaddr": "fa:16:3e:dc:86:b9",
  "ipv4": [
   {
    "address": "10.1.12.176",
    "cidr": "10.1.12.176/20"
   }
  ],
  "ipv6": [
   {
    "address": "fe80::f816:3eff:fedc:86b9",
    "cidr": "fe80::f816:3eff:fedc:86b9/64"
   }
  ],
  "mtu": 1500
}
```

To retrieve an exact value, "walk" the tree all the way to a final element:

```shell
$ elements get -p system.interfaces.ens3.ipv4.0.address

10.1.12.176
```

### Elements Daemon

Elements is able to be run in a daemon mode and accessed over HTTP. In one terminal, run:

```shell
$ elements serve
```

and in another, run:

```shell
$ curl localhost:8888/elements/system/interfaces/ens3
```

## External Elements

Static JSON files and executable files can be placed under `/etc/elements/elements.d`. They will automatically be executed and read when `elements` is run.

The output of these files must be valid JSON. For example, given the executable file `/etc/elements/elements.d/foo.sh`:

```bash
#!/bin/bash

echo '{"hello": "world"}'
```

Elements will output:

```shell
$ elements get
{
  "external": {
   "foo": {
    "hello": "world"
   }
  },
	"system": {
	...
```

If the file is executable (ie: `chmod +x`), Elements will execute it. Non-executable files will be read directly.

## Compile from Source

1. Setup a standard Go environment.
2. Run: `go get github.com/jtopjian/elements/...`
3. Run: `cd $GOROOT/src/github.com/jtopjian/elements`
4. Run: `go build -o elements cmd/*.go`

## History and Credits

Elements was originally a large fork of [Terminus](https://github.com/kelseyhightower/terminus). The move to `gopsutil` was inspired by [go-facter](https://github.com/zstyblik/go-facter).
