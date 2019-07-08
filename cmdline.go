// Copyright Oula Kuuva 2019
// This is pretty much C&P from https://github.com/joneskoo/ruuvi-prometheus/blob/master/cmdline.go
// Thanks Joonas!
package main

import (
	"flag"
)

type settings struct {
	debug          bool
	useDefaults    bool
	configFilePath string
}

func parseSettings() (cmdline settings) {
	flag.BoolVar(&cmdline.debug, "debug", false, "Debug output")
	flag.BoolVar(&cmdline.useDefaults, "builtin-defaults", false, "Use builtin default mapping values")
	flag.StringVar(&cmdline.configFilePath, "config", "", "Path to configuration yaml")
	flag.Parse()
	return cmdline
}
