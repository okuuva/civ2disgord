package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/okuuva/civ2disgord/civ2disgord"
)

func HandleRequest(civ6Message civ2disgord.Civ6Message) error {
	discordMessage, err := civ6Message.NewDefaultDiscordMessage(&civ2disgord.DefaultDiscordConfig, false)
	if err != nil {
		return err
	}
	responses, errs := discordMessage.SendMessage()
	if !checkResponses(responses, nil) {
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