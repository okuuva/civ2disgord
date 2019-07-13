package civ2disgord

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
	"os"
	"strings"
	"testing"
)

var testMessage = fmt.Sprint("{'value1':'RegularGaem', 'value2':'SteamNick2', 'value3':'666'}")
var referenceCivMessage = Civ6Message{
	Value1: "RegularGaem",
	Value2: "SteamNick2",
	Value3: "666",
}

var testContent = fmt.Sprint("Hey <@23123discordIDhere123123>, it's time to take your turn #666 in 'RegularGaem'!")
var testWebhooks = []string{
	"https://when-all-goes-bonkers",
	"https://discordapp.com/webhook0",
}
var referenceDiscordMessage = DiscordMessage{
	Content: testContent,
	webhooks: testWebhooks,
}

func TestParseConfig(t *testing.T) {
	testFilePath := "./config.yml"
	f, err := os.Open(testFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", testFilePath)
		t.Errorf("Error: %s", err)
		return
	}
	config, err := ParseConfig(f)
	if err != nil {
		t.Errorf("Failed to parse %s", testFilePath)
		t.Errorf("Error: %s", err)
		return
	}
	if !cmp.Equal(config, DefaultDiscordConfig) {
		t.Errorf("Parsed config doesn't match the reference!")
		t.Errorf("Reference: %+v", DefaultDiscordConfig)
		t.Errorf("Parsed:    %+v", config)
	}
}

func TestDiscordConfig_DiscordID(t *testing.T) {
	discordID := DefaultDiscordConfig.DiscordID("SteamNick2")
	expected := "23123discordIDhere123123"
	if discordID != expected {
		t.Errorf("Unexpected DiscordID")
		t.Errorf("Expected: %s", expected)
		t.Errorf("Got:      %s", discordID)
	}
}

func TestDiscordConfig_Webhook(t *testing.T) {
	webhook := DefaultDiscordConfig.Webhook("RegularGaem")
	expected := "https://discordapp.com/webhook0"
	if webhook != expected {
		t.Errorf("Unexpected webhook!")
		t.Errorf("Expected: %s", expected)
		t.Errorf("Got:      %s", webhook)
	}
}

func TestParseMessage(t *testing.T) {
	message, err := ParseMessage(strings.NewReader(testMessage))
	if err != nil {
		t.Errorf("Failed to parse test message %s", testMessage)
		t.Errorf("Error: %s", err)
		return
	}
	if !cmp.Equal(message, &referenceCivMessage) {
		t.Errorf("Parsed message doesn't match the reference!")
		t.Errorf("Reference: %+v", &referenceCivMessage)
		t.Errorf("Parsed:    %+v", message)
	}
}

func TestCiv6Message_Game(t *testing.T) {
	game := referenceCivMessage.Game()
	expected := "RegularGaem"
	if game != expected {
		t.Errorf("Unexpected game name!")
		t.Errorf("Expected: %s", expected)
		t.Errorf("Got:      %s", game)
	}
}

func TestCiv6Message_Player(t *testing.T) {
	player := referenceCivMessage.Player()
	expected := "SteamNick2"
	if player != expected {
		t.Errorf("Unexpected player name!")
		t.Errorf("Expected: %s", expected)
		t.Errorf("Got:      %s", player)
	}
}

func TestCiv6Message_TurnNumber(t *testing.T) {
	turn := referenceCivMessage.TurnNumber()
	expected := "666"
	if turn != expected {
		t.Errorf("Unexpected turn number!")
		t.Errorf("Expected: %s", expected)
		t.Errorf("Got:      %s", turn)
	}
}

func TestNewDefaultDiscordMessage(t *testing.T) {
	discordMessage := NewDefaultDiscordMessage("23123discordIDhere123123", "RegularGaem", "666", testWebhooks)
	if !cmp.Equal(&referenceDiscordMessage, discordMessage, cmp.AllowUnexported(DiscordMessage{})) {
		t.Errorf("Generated message doesn't match the reference!")
		t.Errorf("Reference: %+v", &referenceDiscordMessage)
		t.Errorf("Parsed:    %+v", discordMessage)
	}
}

func TestCiv6Message_NewDefaultDiscordMessage(t *testing.T) {
	discordMessage, err := referenceCivMessage.NewDefaultDiscordMessage(&DefaultDiscordConfig, true)
	if err != nil {
		t.Errorf("Failed to generate test Discord message")
		t.Errorf("Error: %s", err)
		return
	}
	if !cmp.Equal(&referenceDiscordMessage, discordMessage, cmp.AllowUnexported(DiscordMessage{})) {
		t.Errorf("Generated message doesn't match the reference!")
		t.Errorf("Reference: %+v", &referenceDiscordMessage)
		t.Errorf("Parsed:    %+v", discordMessage)
	}
}

func TestCiv6Message_NewDefaultDiscordMessageNoWebhooks(t *testing.T) {
	noWebhooksConfig := DiscordConfig{}
	err := copier.Copy(&DefaultDiscordConfig, &noWebhooksConfig)
	if err != nil {
		t.Errorf("Failed to copy DefaultDiscordConfig!")
		t.Errorf("Error: %s", err)
		return
	}
	noWebhooksConfig.Webhooks = map[string]string{}
	_, err = referenceCivMessage.NewDefaultDiscordMessage(&DefaultDiscordConfig, true)
	if err == nil {
		t.Errorf("Generating message from DiscordConfig without webhooks didn't return an error!")
	}
}

func TestCiv6Message_NewDefaultDiscordMessageNoDiscordID(t *testing.T) {
	noMatchingDiscordID := DiscordConfig{}
	err := copier.Copy(&DefaultDiscordConfig, &noMatchingDiscordID)
	if err != nil {
		t.Errorf("Failed to copy DefaultDiscordConfig!")
		t.Errorf("Error: %s", err)
		return
	}
	noMatchingDiscordID.Players = map[string]string{}
	_, err = referenceCivMessage.NewDefaultDiscordMessage(&DefaultDiscordConfig, false)
	if err != nil {
		t.Errorf("Generating message without matching DiscordID while requireDiscordID = false failed!")
	}
}

func TestCiv6Message_NewDefaultDiscordMessageNoDiscordIDWhileRequired(t *testing.T) {
	noMatchingDiscordID := DiscordConfig{}
	err := copier.Copy(&DefaultDiscordConfig, &noMatchingDiscordID)
	if err != nil {
		t.Errorf("Failed to copy DefaultDiscordConfig!")
		t.Errorf("Error: %s", err)
		return
	}
	noMatchingDiscordID.Players = map[string]string{}
	_, err = referenceCivMessage.NewDefaultDiscordMessage(&DefaultDiscordConfig, true)
	if err == nil {
		t.Errorf("Generating message without matching DiscordID while requireDiscordID = true didn't fail!")
	}
}
