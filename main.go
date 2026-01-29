package main

import (
	"errors"
	"log"

	"context"

	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/discord/auth"
	"github.com/dreamervulpi/tourneyBot/internal/discord/bot"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("LoadEnv: %v\n", err)
	}

	ctx := context.Background()

	discordClient, err := auth.NewClient(ctx)
	if err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}

	// 3. Делаем тестовый запрос к GraphQL API
	log.Println("Запрашиваем данные профиля...")
	auth.TestStartGGCall(discordClient)

	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded: ").Error() + err.Error())
	} else {
		tournament, err := config.LoadTournament("config/tournament.toml")
		if err != nil {
			log.Println(errors.New("not loaded: ").Error() + err.Error())
		} else {
			if err := bot.Start(discordClient, cfg, tournament); err != nil {
				log.Println(err.Error())
				// TODO: SAVE LOGS IN TEXT FILE
			}
		}
	}
}
