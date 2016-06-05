# Elements

Get information about a system.

## Install

Download the latest [release](https://github.com/jtopjian/elements/releases).

## Usage

Elements is able to detect a standard set of information about Linux systems. Simply run `elements` on the command line and see the result:

```shell
$ elements get

  "System": {
       "Architecture": "x86_64",
       "BlockDevices": {
        "vda": {
          "Device": "vda",
          "IOTicks": 6547904,
          "InFlight": 0,
          "ReadIOs": 124315,
       ...
```

The output of Elements is JSON to make it easy to import into other utilities and systems.

To retrieve a subset of elements, run:

```shell
$ elements get -p system.BlockDevices.vda

{
  "Device": "vda",
  "IOTicks": 6548052,
  "InFlight": 0,
  "ReadIOs": 124315,
  ...
```

To retrieve an exact value, run:

```shell
$ elements get -p System.BlockDevices.vda.Device

vda
```

### Elements Daemon

Elements is able to be run in a daemon mode and accessed over HTTP. In one terminal, run:

```shell
$ elements serve
```

and in another, run:

```shell
$ curl localhost:8888/elements/System/BlockDevices/vda
{
  "Device": "vda",
  "IOTicks": 6548052,
  "InFlight": 0,
  "ReadIOs": 124315,
  ...
```

## External Elements

Static JSON files and executable files can be placed under `/etc/elements/elements.d`. They will automatically be executed and read when `elements` is run.

The output of these files must be valid JSON. For example:

```bash
#!/bin/bash

echo '{"hello": "world"}'
```

## Compile from Source

1. Setup a standard Go environment.
2. Run: `go get github.com/jtopjian/elements`
3. Run `cd $GOROOT/src/github.com/jtopjian/elements`
4. Run: `go build -o elements cmd/*.go`

## History and Credits

Elements was originally a large fork of [Terminus](https://github.com/kelseyhightower/terminus) which you can find [here](https://github.com/jtopjian/terminus).
