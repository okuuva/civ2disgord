// Copyright Oula Kuuva 2019
package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/jinzhu/copier"
	"github.com/okuuva/civ2disgord/civ2disgord"
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

	var config civ2disgord.DiscordConfig
	if cmdline.configFilePath != "" {
		var f io.Reader
		var err error
		f, err = os.Open(cmdline.configFilePath)
		logger.checkFatal(err, "Could not open config file", 1)
		config, err = civ2disgord.ParseConfig(&f)
		logger.checkFatal(err, "Could not parse config file", 1)
	} else if cmdline.useDefaults {
		err := copier.Copy(&config, &civ2disgord.DefaultDiscordConfig)
		logger.checkFatal(err, "Failed to access default config", 2)
	} else {
		logger.checkFatal(errors.New("no config provided"), "No config provided", 3)
	}
	if len(cmdline.messages) == 0 {
		logger.checkFatal(errors.New("no message given"), "No message given", 3)
	}
	}
}
