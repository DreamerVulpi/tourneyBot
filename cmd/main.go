package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"runtime/debug"

	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/db"
	"github.com/dreamervulpi/tourneyBot/internal/usecase/bot/discord"
	"github.com/joho/godotenv"
)

func initLogger(logDir string) (*os.File, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create log directory: %v", err)
	}

	currentTime := time.Now().Format("02-01-2006_15-04-05")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("tourneyHelper_%s.log", currentTime))
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("--- Logger initialized ---")
	return file, nil
}

func main() {
	logDir := "../logs"
	f, err := initLogger(logDir)
	if err != nil {
		fmt.Printf("Can't launch logging: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Critical error\n Why? %v\n Stack:\n%s", r, debug.Stack())
			fmt.Println("Programm closed with error. More details in folder logs")
		}
	}()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("LoadEnv: %v\n", err)
	}

	// ctx := context.Background()

	// // Challonge
	// chAuth := &auth.AuthClient{
	// 	Config:     auth.GetChallongeOauth2(),
	// 	TokenFile:  "token_challonge.json",
	// 	HTTPClient: &http.Client{},
	// }

	// // auth.TestChallongeCall(chAuth)
	// token, err := chAuth.GetAccessToken("token_challonge.json")
	// if err != nil {
	// 	log.Printf("can't get token Challonge: %v\n", err)
	// }

	// ch := challonge.NewClient(chAuth.HTTPClient, token)
	// tournament, err := ch.GetTournament(ctx, "https://challonge.com/ru/tournamentdciii")
	// if err != nil {
	// 	log.Printf("err | %v", err)
	// }
	// log.Println(tournament.Name)
	// log.Println(tournament.Description)

	// matches, err := ch.GetMatches(ctx, "https://challonge.com/ru/tournamentdciii")
	// log.Println(matches)

	// p, err := ch.GetParticipant(ctx, "https://challonge.com/ru/tournamentdciii", "112579133")
	// if err != nil {
	// 	log.Printf("err | %v", err)
	// }
	// log.Println(p)

	// Для Discord
	dsAuth := &auth.AuthClient{
		Config:    auth.GetDiscordOauth2(),
		TokenFile: "token_discord.json",
	}

	ggAuth := &auth.AuthClient{
		Config:    auth.GetStartggOauth2(),
		TokenFile: "token_startgg.json",
	}

	log.Println("Запрашиваем данные профиля...")

	cfg, err := config.LoadConfig(config.GetAbsPath("config/config.toml"))
	if err != nil {
		log.Println(errors.New("not loaded: ").Error() + err.Error())
	} else {
		pool, err := db.NewPool()
		if err != nil {
			log.Println(err)
			return
		}
		tournament, err := config.LoadTournament(config.GetAbsPath("config/tournament.toml"))
		if err != nil {
			log.Println(errors.New("not loaded: ").Error() + err.Error())
		} else {
			dh := discord.DiscordHandler{}
			if err := dh.Start(dsAuth, ggAuth, pool, cfg, tournament); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
