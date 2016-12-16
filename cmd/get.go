package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jtopjian/elements/lib"
)

var cmdGet cli.Command

func init() {
	cmdGet = cli.Command{
		Name:   "get",
		Usage:  "get one or more elements",
		Action: actionGet,
		Flags: []cli.Flag{
			&flagConfigDir,
			&flagDebug,
			&flagElementPath,
			&flagFormat,
		},
	}
}

func actionGet(c *cli.Context) {
	config := lib.Config{
		Directory:    c.String("configdir"),
		OutputFormat: c.String("format"),
		Path:         c.String("path"),
	}

	elements := lib.Elements{
		Config: config,
	}

	output, err := elements.Get()
	if err != nil {
		errAndExit(err)
	}

	formattedOutput, err := lib.PrintJSON(output)
	if err != nil {
		errAndExit(err)
	}

	fmt.Printf("%s", formattedOutput)
}
