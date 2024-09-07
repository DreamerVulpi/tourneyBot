package bot

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/config"
)

func Start(cfg config.Config) error {
	// TODO: check variables

	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	cmdHandler := commandHandler{
		guildID: cfg.Discord.GuildID,
		stop:    make(chan struct{}),
	}

	commandHandlers["check"] = cmdHandler.check
	commandHandlers["start-sending"] = cmdHandler.start_sending
	commandHandlers["stop-sending"] = cmdHandler.stop_sending
	commandHandlers["set-event"] = cmdHandler.setEvent
	commandHandlers["set-guild-id"] = cmdHandler.setGuildID
	commandHandlers["edit-invite-message"] = cmdHandler.editInviteMessage

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	fmt.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := session.ApplicationCommandCreate(cfg.Discord.AppID, cfg.Discord.GuildID, command)
		if err != nil {
			fmt.Printf("can't create '%v' command: %v\n", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	fmt.Println("the bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("press Ctrl+C to exit")
	<-stop

	fmt.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(cfg.Discord.AppID, cfg.Discord.GuildID, v.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	fmt.Println("gracefully shutting down.")
	return nil
}
