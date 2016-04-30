package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

var flagDebug cli.BoolFlag
var flagConfigDir cli.StringFlag
var flagElementPath cli.StringFlag
var flagListen cli.StringFlag
var defaultConfigDir string

func init() {
	// Figure out the OS
	switch goos := runtime.GOOS; goos {
	case "linux":
		defaultConfigDir = "/etc/elements"
	case "darwin":
		defaultConfigDir = "/etc/elements"
	}

	// flagDebug turns debugging on and off.
	flagDebug = cli.BoolFlag{
		Name:        "debug,d",
		Usage:       "debug mode",
		Destination: &debugMode,
	}

	// flagConfigDir specifies an alternative location to the config directory.
	flagConfigDir = cli.StringFlag{
		Name:  "configdir,c",
		Usage: "Configuration directory",
		Value: defaultConfigDir,
	}

	// flagElementPath specifies the path in the element tree to retrieve.
	flagElementPath = cli.StringFlag{
		Name:  "path,p",
		Usage: "Path in the element tree to retrieve",
	}

	// flagListen specifies the address to listen on for serving elements via http.
	flagListen = cli.StringFlag{
		Name:  "listen,l",
		Usage: "Address to serve elements via http",
		Value: ":8888",
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "element"
	app.Usage = "Obtain facts about a system"
	app.Version = version

	app.Flags = []cli.Flag{
		&flagDebug,
		&flagConfigDir,
		&flagElementPath,
	}

	app.Commands = []cli.Command{
		cmdGet,
		cmdServe,
	}

	app.Run(os.Args)
}
