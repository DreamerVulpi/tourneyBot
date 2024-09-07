package main

import (
	"errors"
	"log"

	"github.com/dreamervulpi/tourneybot/config"
	"github.com/dreamervulpi/tourneybot/internal/bot"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

func main() {
	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded configation"))
	}

	startgg.AuthToken = cfg.Startgg.Token

	bot.Start(cfg)
}
