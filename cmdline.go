// Copyright Oula Kuuva 2019
// This is pretty much C&P from https://github.com/joneskoo/ruuvi-prometheus/blob/master/cmdline.go
// Thanks Joonas!
package main

import (
	"flag"
)

type settings struct {
	debug			bool
	fromEnv			bool
	configFilePath	string
	messages		[]string
}

func parseSettings() (cmdline settings) {
	flag.BoolVar(&cmdline.debug, "debug", false, "Debug output")
	flag.BoolVar(&cmdline.fromEnv, "from-env", false, "Read mappings from environment variables")
	flag.StringVar(&cmdline.configFilePath, "config", "", "Path to configuration yaml")
	flag.Parse()
	cmdline.messages = flag.Args()
	return cmdline
}
