package bot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/config"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

func Start(cfg config.Config, l config.ConfigLobby) error {
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
		guildID:   cfg.Discord.GuildID,
		client:    client,
		stop:      make(chan struct{}),
		dataLobby: l,
		RulesMatches: config.RulesMatches{
			Format:        l.Rules.Format,
			Rounds:        l.Rules.Rounds,
			Duration:      l.Rules.Duration,
			Crossplatform: l.Rules.Crossplatform,
			Stage:         l.Rules.Stage,
		},
		StreamLobby: config.StreamLobby{
			Area:          l.Stream.Area,
			Language:      l.Stream.Language,
			Crossplatform: l.Stream.Crossplatform,
			Conn:          l.Stream.Conn,
			Passcode:      l.Stream.Passcode,
		},
	}

	commandHandlers["check"] = cmdHandler.check
	commandHandlers["start-sending"] = cmdHandler.start_sending
	commandHandlers["stop-sending"] = cmdHandler.stop_sending
	commandHandlers["set-event"] = cmdHandler.setEvent
	commandHandlers["edit-rules"] = cmdHandler.editRuleMatches
	commandHandlers["edit-stream-lobby"] = cmdHandler.editStreamLobby

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	log.Println("adding commands...")
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
