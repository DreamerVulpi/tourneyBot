package main

import (
	"errors"
	"log"

	"github.com/dreamervulpi/tourneybot/config"
	"github.com/dreamervulpi/tourneybot/internal/bot"
)

func main() {
	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded configation: ").Error() + err.Error())
	} else {
		bot.Start(cfg)
	}
}
