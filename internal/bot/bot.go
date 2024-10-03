package bot

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type contactData struct {
	discord string
	gameID  string
}

// Get discord contacts from csv file
func loadCSV(nameFile string) map[string]contactData {
	contacts := map[string]contactData{}
	f, err := os.Open("config/" + nameFile)
	if err != nil {
		log.Println(err)
		return map[string]contactData{}
	} else {
		if len(nameFile) != 0 {
			defer f.Close()

			csvReader := csv.NewReader(f)
			records, _ := csvReader.ReadAll()

			// Search index for get data
			var indexDiscordColumn int
			var indexGamerTagColumn int
			var indexConnectColumn int
			for index, column := range records[0] {
				parts := strings.SplitN(column, " ", -1)
				for _, part := range parts {
					if part == "Discord!" {
						indexDiscordColumn = index
					}
				}
				if column == "Short GamerTag" {
					indexGamerTagColumn = index
				}
				if column == "Connect" {
					indexConnectColumn = index
				}
			}

			for i, attendee := range records {
				if i == 0 {
					continue
				}

				var discordID string
				if len(attendee[indexDiscordColumn]) != 0 {
					discordID = attendee[indexDiscordColumn]
				} else {
					discordID = "N/D"
				}

				var gameID string
				if len(attendee[indexConnectColumn]) != 0 {
					rawTekkenID := strings.SplitN(attendee[indexConnectColumn], " ", -1)
					gameID = strings.ReplaceAll(rawTekkenID[1], ",", "")
				} else {
					gameID = "N/D"
				}

				contacts[attendee[indexGamerTagColumn]] = contactData{
					discord: discordID,
					gameID:  gameID,
				}
			}
		}
	}
	return contacts
}

func Start(cfg config.Config, t config.ConfigTournament) error {
	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	client := startgg.NewClient(cfg.Startgg.Token, &http.Client{
		Timeout: time.Second * 10,
	})

	cmdHandler := commandHandler{
		guildID:    cfg.Discord.GuildID,
		client:     client,
		stop:       make(chan struct{}),
		tournament: t,
		rulesMatches: config.RulesMatches{
			Format:        t.Rules.Format,
			Rounds:        t.Rules.Rounds,
			Duration:      t.Rules.Duration,
			Crossplatform: t.Rules.Crossplatform,
			Stage:         t.Rules.Stage,
		},
		streamLobby: config.StreamLobby{
			Area:          t.Stream.Area,
			Language:      t.Stream.Language,
			Crossplatform: t.Stream.Crossplatform,
			Conn:          t.Stream.Conn,
			Passcode:      t.Stream.Passcode,
		},
		logo:            "https://i.imgur.com/n9SG5IL.png",
		logoTournament:  t.Logo.Img,
		appID:           cfg.Discord.AppID,
		rolesIdList:     cfg.Roles,
		discordContacts: loadCSV(t.Csv.NameFile),
	}

	commandHandlers["check"] = cmdHandler.viewData
	commandHandlers["start-sending"] = cmdHandler.start_sending
	commandHandlers["stop-sending"] = cmdHandler.stop_sending
	commandHandlers["set-event"] = cmdHandler.setEvent
	commandHandlers["edit-rules"] = cmdHandler.editRuleMatches
	commandHandlers["edit-stream-lobby"] = cmdHandler.editStreamLobby
	commandHandlers["edit-logo-tournament"] = cmdHandler.editLogoTournament
	commandHandlers["contacts"] = cmdHandler.viewContacts

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	log.Println("adding commands...")
	commands := cmdHandler.commands()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := session.ApplicationCommandCreate(cfg.Discord.AppID, cfg.Discord.GuildID, command)
		log.Printf("%v\n", command.Name)
		if err != nil {
			log.Printf("can't create '%v' command: %v\n", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	log.Println("the bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press Ctrl+C to exit")
	<-stop

	log.Println("removing commands...")
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(cfg.Discord.AppID, cfg.Discord.GuildID, v.ID)
		log.Printf("%v\n", v.Name)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
	log.Println("gracefully shutting down!")
	return nil
}
