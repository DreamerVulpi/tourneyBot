package main

import (
	"errors"
	"log"

	"context"

	"github.com/dreamervulpi/tourneyBot/config"
	auth "github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/db"
	"github.com/dreamervulpi/tourneyBot/internal/discord/bot"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("LoadEnv: %v\n", err)
	}

	ctx := context.Background()
	// Для Discord
	dsAuth := &auth.AuthClient{
		Config:    auth.GetDiscordOauth2(),
		TokenFile: "token_discord.json",
	}
	if err := dsAuth.Init(ctx); err != nil {
		log.Fatal(err)
	}

	// Для Start.gg
	stAuth := &auth.AuthClient{
		Config:    auth.GetStartggOauth2(),
		TokenFile: "token_startgg.json",
	}
	if err := stAuth.Init(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Запрашиваем данные профиля...")
	auth.TestStartGGCall(stAuth)

	cfg, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Println(errors.New("not loaded: ").Error() + err.Error())
	} else {
		pool, err := db.NewPool()
		if err != nil {
			log.Println(err)
			return
		}
		tournament, err := config.LoadTournament("config/tournament.toml")
		if err != nil {
			log.Println(errors.New("not loaded: ").Error() + err.Error())
		} else {
			if err := bot.Start(stAuth.HTTPClient, dsAuth, pool, cfg, tournament); err != nil {
				log.Println(err.Error())
				// TODO: SAVE LOGS IN TEXT FILE
			}
		}
	}
}
