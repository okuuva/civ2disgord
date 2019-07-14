package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/okuuva/civ2disgord/civ2disgord"
	"io/ioutil"
	"os"
)

func HandleRequest(civ6Message civ2disgord.Civ6Message) error {
	logger := newLogger(ioutil.Discard, os.Stdout, os.Stderr)
	discordMessage, err := civ6Message.NewDefaultDiscordMessage(&civ2disgord.DefaultDiscordConfig, false)
	if err != nil {
		return err
	}
	responses, errs := discordMessage.SendMessage()
	if !checkResponses(responses, logger) {
		return fmt.Errorf("failed to send message")
	}
	err = checkErrors(errs)
	if err != nil {
		return fmt.Errorf("failed to send message")
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}