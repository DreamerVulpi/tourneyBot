package bot

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	AuthToken             string
	GuildID               string
	AppID                 string
	Slug                  string
	TemplateInviteMessage string
)

func SetAuthToken(token string) {
	AuthToken = token
}

func SetGuildID(id string) {
	GuildID = id
}

func SetAppID(id string) {
	AppID = id
}

func SetTemplateInviteMessage(template string) {
	TemplateInviteMessage = template
}

func SetSlug(slug string) {
	Slug = slug
}

func slug() bool {
	return len(Slug) > 0
}

func server() bool {
	return len(GuildID) > 0
}

func app() bool {
	return len(AppID) > 0
}

func token() bool {
	return len(AuthToken) > 0
}

func Start() error {
	if !app() {
		return errors.New("appID is empty")
	}

	if !server() {
		return errors.New("serverID(guildID) is empty")
	}

	if !token() {
		return errors.New("authToken is empty")
	}

	session, err := discordgo.New(AuthToken)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(AppID, GuildID, command)
		if err != nil {
			fmt.Printf("can't create '%v' command: %v\n", command.Name, err)
		}
	}

	defer session.Close()

	fmt.Println("the bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("press Ctrl+C to exit")
	<-stop

	// TODO: Add functional for deleting commands

	fmt.Println("gracefully shutting down.")
	return nil
}
