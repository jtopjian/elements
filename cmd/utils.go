package main

import (
	"fmt"
	"log"
	"os"
)

var debugMode bool
var debug debugger

type debugger bool

// debugger.Printf allows printing logs only when debug is true.
func (d debugger) Printf(format string, args ...interface{}) {
	if debugMode {
		log.Printf(format, args...)
	}
}

// errAndExit is a convenience function to print an error and exit.
func errAndExit(e error) {
	fmt.Printf("Error: %s\n", e)
	os.Exit(1)
}
