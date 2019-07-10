// Copyright Oula Kuuva 2019
package main

import (
	"io"
	"io/ioutil"
	"os"
)

const (
	usage = "Halp"
)

func main() {
	cmdline := parseSettings()
	var debugPipe io.Writer
	if cmdline.debug {
		debugPipe = os.Stdout
	} else {
		debugPipe = ioutil.Discard
	}
	logger := newLogger(debugPipe, os.Stdout, os.Stderr)

	var config *config
	var err error
	if cmdline.configFilePath != "" {
		var f io.Reader
		f, err = os.Open(cmdline.configFilePath)
		logger.checkFatal(err, "Could not open config file", 1)
		config, err = parseConfig(&f)
		logger.checkFatal(err, "Could not parse config file", 2)
	}
}
