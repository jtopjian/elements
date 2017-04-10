# Elements

Elements retrieves information about a system such as CPU/Processors, disks, memory, network interfaces, and cloud metadata. It is meant to be used to query a system about its static attributes.

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

### Output Formats

Elements can be printed out in two different formats: JSON and shell. By default, JSON will be used, but shell is useful for within shell scripts:

```bash
#!/bin/bash

eval $(elements get -f shell)
echo $elements_system_interfaces_ens3_ipv4_0_address
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

You may specify the output format via the `?format=` query parameter, like so:

```shell
$ curl localhost:8888/elements/system/interfaces/ens3?format=shell
```

## External Elements in `elements.d/`

You may extend the facts elements reports by putting external sources with valid
JSON in a configurable `elements.d` directory. By default, this directory is
`/etc/elements/elements.d`.

Executables and static JSON files placed in `/etc/elements/elements.d` will automatically be executed and/or read when `elements` is run. If the file is executable (ie: `chmod +x`), Elements will execute it. Non-executable files with a `.json` extension will be read directly. Files in the `elements.d` directory that do not match these criteria will be ignored. Files and executables that do not contain or produce valid JSON will be ignored.
 
For example, given the executable file `/etc/elements/elements.d/foo.sh`:

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

You my specify an alternate directory location for your `elements.d` directory
with the `-c | --configdir` flag. For example, to source external elements from
the `/var/lib/cloud/data/elements.d` directory, you would specify:

```shell
$ elements get -c /var/lib/cloud/data
```

## Compile from Source

1. Setup a standard Go environment.
2. Run: `go get github.com/jtopjian/elements/...`
3. Run: `cd $GOROOT/src/github.com/jtopjian/elements`
4. Run: `go build -o elements cmd/*.go`

## History and Credits

Elements was originally a large fork of [Terminus](https://github.com/kelseyhightower/terminus). The move to `gopsutil` was inspired by [go-facter](https://github.com/zstyblik/go-facter).
