package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/config"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type commandHandler struct {
	slug    string
	guildID string
	// FIXME: Change to more args (command editTemplate)
	templateInviteMessage string
	// FIXME: Change to more args (command editTemplate)
	templateStreamMessage string
	stop                  chan struct{}
	m                     *discordgo.MessageCreate
	client                *startgg.Client
	dataLobby             config.ConfigLobby
}

func response(s *discordgo.Session, i *discordgo.InteractionCreate, text string) error {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: text,
			},
		},
	)
	if err != nil {
		return errors.New("can't respond on message")
	}
	return nil
}

func responseSetted(s *discordgo.Session, i *discordgo.InteractionCreate, msgformat string, margs []interface{}) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
	if err != nil {
		return errors.New("can't respond on message")
	}
	return nil
}

// TODO: Add Embed-view
func (cmd *commandHandler) check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	text := fmt.Sprint("server(Guild) ID: " + cmd.guildID + "\n" + "slug: " + cmd.slug + "\n" + "templateInviteMessage: \n" + cmd.templateInviteMessage)
	if err := response(s, i, text); err != nil {
		log.Println(err.Error())
	}
}
func (cmd *commandHandler) start_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := response(s, i, "Start sending..."); err != nil {
		log.Println(err.Error())
	}

	go cmd.SendingMessages(s)
}
func (cmd *commandHandler) stop_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	go func() {
		response(s, i, "stopping...")
	}()

	// Send signal to stop process
	cmd.stop <- struct{}{}

	s.ChannelMessageSend(i.ChannelID, "Stopped!")
}

func (cmd *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()
	cmd.slug = input

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> SLUG: %s\n"

	if err := responseSetted(s, i, msgformat, margs); err != nil {
		log.Println(err.Error())
	}
}

// FIXME: Add more args
func (cmd *commandHandler) editInviteMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> Template: %s\n"
	cmd.templateInviteMessage = input

	if err := responseSetted(s, i, msgformat, margs); err != nil {
		log.Println(err.Error())
	}
}

// TODO: Add new command: urlLogo
// TODO: Add new command: editStreamMessage (with more args)
// TODO: Add new command: Help
// TODO: Add new command: About
