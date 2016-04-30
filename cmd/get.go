package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	elements "github.com/jtopjian/elements/lib"
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
		},
	}
}

func actionGet(c *cli.Context) {
	e, err := elements.New(c.String("configdir"), c.String("path"))
	if err != nil {
		errAndExit(err)
	}

	elements, err := e.Elements2JSON()
	if err != nil {
		errAndExit(err)
	}

	fmt.Printf("%s", elements)

}
