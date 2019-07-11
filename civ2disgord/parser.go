// Package civ2disgord makes Civ 6 Play by Cloud notifications and Discord play ball
package civ2disgord

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// I really do not want to know how the actual game code looks like when the API has this kind of variable naming...
type Civ6Message struct {
	Value1 string // Player name
	Value2 string // Game name
	Value3 string // Turn number
}

// Add sane naming for entries without a need for a type conversion
func (civMessage *Civ6Message) Player() string     { return civMessage.Value1 }
func (civMessage *Civ6Message) Game() string       { return civMessage.Value2 }
func (civMessage *Civ6Message) TurnNumber() string { return civMessage.Value3 }

func (civMessage *Civ6Message) NewDefaultDiscordMessage(config *DiscordConfig, requireDiscordID bool) (*DiscordMessage, error) {
	player := civMessage.Player()
	game := civMessage.Game()
	turn := civMessage.TurnNumber()
	discordID := config.DiscordID(player)
	if discordID == "" {
		if requireDiscordID {
			return nil, fmt.Errorf("Could not find DiscordID for player %s", player)
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
		err = fmt.Errorf("Could not find webhook for game %s", game)
	} else {
		webhooks = append(webhooks, webhook)
	}
	discordMessage := NewDefaultDiscordMessage(discordID, game, turn, &webhooks)
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
	Content  string
	webhooks []string
}

func NewDefaultDiscordMessage(player, game, turn string, webhooks *[]string) *DiscordMessage {
	var discordMessage DiscordMessage
	discordMessage.Content = fmt.Sprintf("Hey <@%s>, it's time to take your turn #%s in '%s'!", player, game, turn)
	discordMessage.webhooks = webhooks
	return &discordMessage
}

type DiscordConfig struct {
	Players      map[string]string
	Webhooks     map[string]string
	DebugWebhook string
}

func (config *DiscordConfig) DiscordID(player string) string {return config.Players[player]}
func (config *DiscordConfig) Webhook(game string) string {return config.Webhooks[game]}

func ParseConfig(configFile *io.Reader) (*DiscordConfig, error) {
	decoder := yaml.NewDecoder(*configFile)
	var config DiscordConfig
	err := decoder.Decode(&config)
	return &config, err
}
