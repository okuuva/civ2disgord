// Copyright Oula Kuuva 2019
package main

import (
	"io"
	"io/ioutil"
	"os"
)

func check(err error, message string, returnCode int) {
	if err != nil {
		error.Println(message)
		error.Println(err)
		os.Exit(returnCode)
	}
}

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
	var config *config
	var err error
	if cmdline.configFilePath != "" {
		var f io.Reader
		f, err = os.Open(cmdline.configFilePath)
		check(err, "Could not open config file", 2)
		config, err = parseConfig(&f)
		check(err, "Could not parse config file", 3)
	}
}
