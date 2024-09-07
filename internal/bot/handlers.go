package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type commandHandler struct {
	slug                  string
	guildID               string
	templateInviteMessage string
	stop                  chan struct{}
	m                     *discordgo.MessageCreate
	client                *startgg.Client
}

func (cmd *commandHandler) check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond := fmt.Sprint("server(Guild) ID: " + cmd.guildID + "\n" + "slug: " + cmd.slug + "\n" + "templateInviteMessage: \n" + cmd.templateInviteMessage)
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

	go cmd.SendingMessages(s)

	if err != nil {
		log.Println(errors.New("can't respond on message"))
	}
}
func (cmd *commandHandler) stop_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	go func() {
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Stopping...",
				},
			},
		)

		if err != nil {
			log.Println(errors.New("can't respond on message"))
		}
	}()

	// Send signal to stop sending messages
	cmd.stop <- struct{}{}

	s.ChannelMessageSend(i.ChannelID, "Stopped!")

}
func (cmd *commandHandler) setGuildID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	cmd.guildID = input

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> GuildID: %s\n"

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
	cmd.slug = input

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
	cmd.templateInviteMessage = input

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
