// Package civ2disgord makes Civ 6 Play by Cloud notifications and Discord play ball
package civ2disgord

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func getEnv(key string) string {
	key = b64.RawURLEncoding.EncodeToString([]byte(key))
	return os.Getenv(key)
}

// I really do not want to know how the actual game code looks like when the API has this kind of variable naming...
type Civ6Message struct {
	Value1 string // Game name
	Value2 string // Player name
	Value3 string // Turn number
}

// Add sane naming for entries without a need for a type conversion
func (civMessage *Civ6Message) Game() string       { return civMessage.Value1 }
func (civMessage *Civ6Message) Player() string     { return civMessage.Value2 }
func (civMessage *Civ6Message) TurnNumber() string { return civMessage.Value3 }

func (civMessage *Civ6Message) NewDefaultDiscordMessage(config *DiscordConfig, requireDiscordID bool) (*DiscordMessage, error) {
	player := civMessage.Player()
	game := civMessage.Game()
	turn := civMessage.TurnNumber()
	discordID := config.DiscordID(player)
	if discordID == "" {
		if requireDiscordID {
			return nil, fmt.Errorf("could not find DiscordID for player %s", player)
		}
		discordID = player
	}
	var webhooks []string
	webhook := config.DebugWebhook
	if webhook != "" {
		webhooks = append(webhooks, webhook)
	}
	webhook = config.Webhook(game)
	var err error
	if webhook == "" {
		err = fmt.Errorf("could not find webhook for game %s", game)
	} else {
		webhooks = append(webhooks, webhook)
	}
	discordMessage := NewDefaultDiscordMessage(discordID, game, turn, webhooks)
	return discordMessage, err
}

func (civMessage *Civ6Message) NewDefaultDiscordMessageFromEnv(requireDiscordID bool) (*DiscordMessage, error) {
	player := civMessage.Player()
	game := civMessage.Game()
	turn := civMessage.TurnNumber()
	discordID := getEnv(player)
	if discordID == "" {
		if requireDiscordID {
			return nil, fmt.Errorf("could not find DiscordID for player %s", player)
		}
		discordID = player
	}
	debugWebhooks := []string{
		getEnv("global_debug_webhook"),
		getEnv(fmt.Sprintf("%s_debug", game)),
	}
	var webhooks []string
	for _, webhook := range debugWebhooks {
		if webhook != "" {
			webhooks = append(webhooks, webhook)
		}
	}
	var err error
	webhook := getEnv(game)
	if webhook == "" {
		err = fmt.Errorf("could not find webhook for game %s", game)
	} else {
		webhooks = append(webhooks, webhook)
	}
	discordMessage := NewDefaultDiscordMessage(discordID, game, turn, webhooks)
	return discordMessage, err
}

func ParseMessage(messageBody io.Reader) (*Civ6Message, error) {
	// Since yaml is superset of json, we can decode it with yaml decoder
	// Just use strict decoder that handles the few corners between the two
	decoder := yaml.NewDecoder(messageBody)
	decoder.SetStrict(true)
	var message Civ6Message
	err := decoder.Decode(&message)
	return &message, err
}

type DiscordMessage struct {
	Content  string		`json:"content"`
	webhooks []string
}

func (discordMessage *DiscordMessage) SendMessage() (responses []*http.Response, errs []error) {
	for _, webhook := range discordMessage.webhooks {
		resp, err := discordMessage.sendMessageTo(webhook)
		responses = append(responses, resp)
		errs = append(errs, err)
	}
	return responses, errs
}

func (discordMessage *DiscordMessage) sendMessageTo(url string) (*http.Response, error) {
	payload, err := json.Marshal(discordMessage)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(payload))
}

func NewDefaultDiscordMessage(player, game, turn string, webhooks []string) *DiscordMessage {
	var discordMessage DiscordMessage
	discordMessage.Content = fmt.Sprintf("Hey <@%s>, it's time to take your turn #%s in '%s'!", player, turn, game)
	discordMessage.webhooks = webhooks
	return &discordMessage
}

type DiscordConfig struct {
	Players      map[string]string	`yaml:"players"`
	Webhooks     map[string]string	`yaml:"webhooks"`
	DebugWebhook string				`yaml:"debug-webhook"`
}

func (config *DiscordConfig) DiscordID(player string) string {return config.Players[player]}
func (config *DiscordConfig) Webhook(game string) string {return config.Webhooks[game]}

func ParseConfig(configFile io.Reader) (DiscordConfig, error) {
	decoder := yaml.NewDecoder(configFile)
	var config DiscordConfig
	err := decoder.Decode(&config)
	return config, err
}
