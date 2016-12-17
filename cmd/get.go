package main

import (
	"fmt"

	"github.com/codegangsta/cli"

	e "github.com/jtopjian/elements/lib/elements"
	o "github.com/jtopjian/elements/lib/output"
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
	eConfig := e.Config{
		Directory: c.String("configdir"),
		Path:      c.String("path"),
	}

	elements := e.Elements{
		Config: eConfig,
	}

	oConfig := o.Config{
		Format: c.String("format"),
	}

	output := o.Output{
		Config: oConfig,
	}

	collectedElements, err := elements.Get()
	if err != nil {
		errAndExit(err)
	}

	formattedOutput, err := output.Generate(collectedElements)
	if err != nil {
		errAndExit(err)
	}

	fmt.Printf("%s", formattedOutput)
}
