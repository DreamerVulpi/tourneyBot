package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type commandHandler struct {
	stop chan struct{}
}

func (cmd *commandHandler) check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond := fmt.Sprint("server(Guild) ID: " + GuildID + "\n" + "slug: " + Slug + "\n" + "templateInviteMessage: \n" + TemplateInviteMessage)
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
func (cmd *commandHandler) start_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	go SendingMessages(s, m, cmd.stop)

	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
func (cmd *commandHandler) stop_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Printf("1111111111")
	cmd.stop <- struct{}{}
	fmt.Printf("2222222222")
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Stopped",
			},
		},
	)

	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
func (cmd *commandHandler) setGuildID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> GuildID: %s\n"
	SetGuildID(input)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}
func (cmd *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> SLUG: %s\n"
	SetSlug(input)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}
func (cmd *commandHandler) editInviteMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> Template: %s\n"
	SetTemplateInviteMessage(input)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}
