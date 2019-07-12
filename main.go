// Copyright Oula Kuuva 2019
package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	var responses []*http.Response
	var errs []error
	for _, message := range cmdline.messages {
		civMessage, err := civ2disgord.ParseMessage(strings.NewReader(message))
		logger.checkFatal(err, "Failed to parse message", 4)
		discordMessage, err := civMessage.NewDefaultDiscordMessage(&config, false)
		logger.checkFatal(err, "Failed to construct DiscordMessage", 5)
		responses, errs = discordMessage.SendMessage()
	}
	for _, err := range errs {
		if err != nil {
			logger.error.Println(err)
			os.Exit(5)
		}
	}
	var fail bool
	for _, response := range responses {
		logger.debug.Println(response)
		logger.debug.Println(response.Request)
		url := response.Request.URL.String()
		if response.StatusCode != 204 {
			logger.error.Printf("Failed to send message to %s!", url)
			fail = true
		} else {
			logger.info.Printf("Successfully sent meggase to %s", url)
		}
	}
	if fail {
		os.Exit(5)
	}
}
