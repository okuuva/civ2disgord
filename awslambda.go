package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/okuuva/civ2disgord/civ2disgord"
	"os"
)

func HandleRequest(civ6Message civ2disgord.Civ6Message) error {
	logger := newLogger(os.Stdout, os.Stdout, os.Stderr)
	logger.debug.Printf("Received message: %+v", civ6Message)
	discordMessage, err := civ6Message.NewDefaultDiscordMessageFromEnv(false, true)
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