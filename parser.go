package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

// I really do not want to know how the actual game code looks like when the API has this kind of variable naming...
type civ6message struct {
	Value1 string // Player name
	Value2 string // Game name
	Value3 string // Turn number
}

// Add sane naming for entries without a need for a type conversion
func (c civ6message) Player() string     { return c.Value1 }
func (c civ6message) Game() string       { return c.Value2 }
func (c civ6message) TurnNumber() string { return c.Value3 }

type config struct {
	Players      map[string]string
	Webhooks     map[string]string
	DebugWebhook string
}

func parseMessage(messageBody *io.Reader) (*civ6message, error) {
	// Since yaml is superset of json, we can decode it with yaml decoder
	// Just use strict decoder that handles the few corners between the two
	decoder := yaml.NewDecoder(*messageBody)
	decoder.SetStrict(true)
	var message civ6message
	err := decoder.Decode(&message)
	return &message, err
}

func parseConfig(configFile *io.Reader) (*config, error) {
	decoder := yaml.NewDecoder(*configFile)
	var config config
	err := decoder.Decode(&config)
	return &config, err
}
