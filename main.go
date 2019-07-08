// Copyright Oula Kuuva 2019
package main

import (
	"io"
	"io/ioutil"
	"os"
)

func main() {
	cmdline := parseSettings()
	var debugPipe io.Writer
	if cmdline.debug {
		debugPipe = os.Stdout
	} else {
		debugPipe = ioutil.Discard
	}
	initLoggers(debugPipe, os.Stdout, os.Stdout, os.Stderr)

	if !cmdline.useDefaults && cmdline.configFilePath == "" {
		error.Println("No -builtin-defaults set nor config path given!")
		os.Exit(1)
	}
}
