// Package civ2disgord makes Civ 6 Play by Cloud notifications and Discord play ball
package civ2disgord

import (
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
