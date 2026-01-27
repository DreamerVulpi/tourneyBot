package main

import (
	"errors"
	"log"

	"context"

	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("LoadEnv: %v\n", err)
		return
	}

	ctx := context.Background()

	client, err := auth.NewClient(ctx)
	if err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}

	// 3. Делаем тестовый запрос к GraphQL API
	log.Println("Запрашиваем данные профиля...")
	auth.TestStartGGCall(client)

	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded: ").Error() + err.Error())
	} else {
		tournament, err := config.LoadTournament("config/tournament.toml")
		if err != nil {
			log.Println(errors.New("not loaded: ").Error() + err.Error())
		} else {
			if err := bot.Start(cfg, tournament); err != nil {
				log.Println(err.Error())
				// TODO: SAVE LOGS IN TEXT FILE
			}
		}
	}
}
