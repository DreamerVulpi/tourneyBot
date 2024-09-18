package main

import (
	"errors"
	"log"

	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/bot"
)

func main() {
	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded: ").Error() + err.Error())
	} else {
		tournament, err := config.LoadTournament("config/tournament.toml")
		if err != nil {
			log.Println(errors.New("not loaded: ").Error() + err.Error())
		} else {
			bot.Start(cfg, tournament)
		}
	}
}
