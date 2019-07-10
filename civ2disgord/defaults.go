// Package civ2disgord makes Civ6 and Discord play ball
// Copyright Oula Kuuva 2019
package civ2disgord

// DefaultPlayers maps Steam nicks into Discord IDs
var DefaultPlayers = map[string]string{
	"SteamNick1": "13123discordIDhere123123",
	"SteamNick2": "23123discordIDhere123123",
}

// DefaultWebhooks maps Civ6 game names to Discord Webhooks
var DefaultWebhooks = map[string]string{
	"RegularGaem":     "https://discordapp.com/webhook0",
	"SupaAwesomeGaem": "https://discordapp.com/webhook1",
}

// DefaultDebugWebhook tells where to yell if something goes bonkers
var DefaultDebugWebhook = "https://when-all-goes-bonkers"
