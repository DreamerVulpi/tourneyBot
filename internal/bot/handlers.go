package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO: Refactor code

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"check": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"start-sending": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Start sending...",
					},
				},
			)

			// FIXME: SEND_MESSAGES
			var m *discordgo.MessageCreate

			SendProcess(s, m)

			if err != nil {
				log.Println(errors.New("can't respond on message"))
			}
		},
		"stop-sending": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "WIP",
					},
				},
			)

			// TODO: Method to correctly stop SendProcess

			if err != nil {
				log.Println(errors.New("can't respond on message"))
			}
		},
		"set-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			fmt.Println(optionMap)

			margs := make([]interface{}, 0, len(options))
			msgformat := ""

			if option, ok := optionMap["slug"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> SLUG: %s\n"
				SetSlug(option.StringValue())
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"set-guild-id": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			fmt.Println(optionMap)

			margs := make([]interface{}, 0, len(options))
			msgformat := ""

			if option, ok := optionMap["guildID"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> GuildID: %s\n"
				SetGuildID(option.StringValue())
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"edit-invite-message": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			msgformat := ""

			if option, ok := optionMap["template"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Input: %s\n"
				SetTemplateInviteMessage(option.StringValue())
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
	}
)
