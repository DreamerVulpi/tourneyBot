package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Start() error {
	session, err := discordgo.New("Bot -")
	if err != nil {
		return err
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}
	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		return err
	}

	defer session.Close()

	fmt.Println("the bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	return nil
}
