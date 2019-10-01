// Copyright Oula Kuuva 2019
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/okuuva/civ2disgord/civ2disgord"
)

const (
	// TODO: create a proper help string
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
		logger.debug.Printf("Reading config from %s", cmdline.configFilePath)
		var f io.Reader
		var err error
		f, err = os.Open(cmdline.configFilePath)
		logger.checkFatal(err, "Could not open config file", 1)
		config, err = civ2disgord.ParseConfig(f)
		logger.checkFatal(err, "Could not parse config file", 1)
		logger.debug.Println("Successfully loaded config")
		if cmdline.toEnv {
			envConfig := config.ToEnvVariables()
			for variable, value := range *envConfig {
				fmt.Printf("%s=%s\n", variable, value)
			}
			os.Exit(0)
		}
	} else if cmdline.fromEnv {
		logger.debug.Println("Reading mapping values from environment variables")
		if err := godotenv.Load(); err != nil {
			logger.debug.Println("Failed to load default environment variables from .env")
		}
		if cmdline.toEnv {
			logger.checkFatal(errors.New("mixed messages"), "No point converting default variables", 3)
		}
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
		var discordMessage *civ2disgord.DiscordMessage
		if cmdline.fromEnv {
			discordMessage, err = civMessage.NewDefaultDiscordMessageFromEnv(false)
		} else {
			discordMessage, err = civMessage.NewDefaultDiscordMessage(&config, false)
		}
		logger.checkFatal(err, "Failed to construct DiscordMessage", 5)
		responses, errs = discordMessage.SendMessage()
	}
	err := checkErrors(errs)
	if err != nil {
		logger.checkFatal(err, "Failed to send message", 5)
	}
	if !checkResponses(responses, logger) {
		os.Exit(5)
	}
}
