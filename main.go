package main

import (
	"errors"
	"log"

	"github.com/dreamervulpi/tourneybot/internal/bot"
	"github.com/dreamervulpi/tourneybot/internal/config"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

func main() {
	cfg, err := config.LoadConfig("internal/config/config.toml")
	if err != nil {
		log.Println(errors.New("Not loaded configation"))
	}

	startgg.AuthToken = cfg.Startgg.Token
	bot.AuthToken = cfg.Discord.Token

	bot.SetAuthToken(cfg.Discord.Token)
	bot.SetServerID(cfg.Discord.GuildID)
	bot.SetAppID(cfg.Discord.AppID)
	bot.Start()
}
