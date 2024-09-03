package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handlerCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "check":
		check(s, i)
	case "start-sending":
		startSending(s, i)
	case "stop-sending":
		stopSending(s, i)
	}
}

func check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond := fmt.Sprint("server(Guild) ID: " + GuildID + "\n" + "slug: " + Slug + "\n" + "templateInviteMessage: \n" + templateInviteMessage)

	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: respond,
			},
		},
	)
	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
func startSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Start sending...",
			},
		},
	)

	var m *discordgo.MessageCreate

	SendProcess(s, m)

	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
func stopSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "WIP",
			},
		},
	)
	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
