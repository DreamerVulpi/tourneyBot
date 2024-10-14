package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type commandHandler struct {
	slug                 string
	guildID              string
	appID                string
	logo                 string
	logoTournament       string
	nameGame             string
	stop                 chan struct{}
	m                    *discordgo.MessageCreate
	client               *startgg.Client
	tournament           config.ConfigTournament
	rulesMatches         config.RulesMatches
	streamLobby          config.StreamLobby
	rolesIdList          config.ConfigRolesIdDiscord
	discordContacts      map[string]contactData
	embedDiscordContacts []*discordgo.MessageEmbed
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

func (cmd *commandHandler) messageEmbed(title string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: title,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: cmd.logo,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cmd.logoTournament,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "by DreamerVulpi | https://www.twitch.tv/dreamervulpi",
			IconURL: "https://i.imgur.com/FcuAfRw.png",
		},
	}
}

func (cmd *commandHandler) responseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed []*discordgo.MessageEmbed) error {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embed,
			},
		},
	)
	return err
}

func (cmd *commandHandler) viewData(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := []*discordgo.MessageEmbed{}

	embed = append(embed, cmd.msgViewData(i.Locale.String()))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("check: can't respond on message"))
	}
}

func (cmd *commandHandler) start_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := cmd.msgResponse(i.Locale.String())

	if cmd.guildID != "" && cmd.slug != "" {
		if err := response(s, i, local.responseMsg.Starting); err != nil {
			log.Println(err.Error())
		}
		go cmd.Process(s)
	} else {
		embed := []*discordgo.MessageEmbed{}

		embed = append(embed, cmd.messageEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "", Value: local.errorMsg.Input},
		}))

		if err := cmd.responseEmbed(s, i, embed); err != nil {
			log.Println(fmt.Errorf("editLogoTournament: %v", local.errorMsg.Respond))
		}
	}
}

func (cmd *commandHandler) stop_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := cmd.msgResponse(i.Locale.String())

	go func() {
		response(s, i, local.responseMsg.Stopping)
	}()

	// Send signal to stop process
	cmd.stop <- struct{}{}

	s.ChannelMessageSend(i.ChannelID, local.responseMsg.Stopped)
}

func (cmd *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	u, err := url.Parse(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		log.Println(err)
	}

	arg := strings.SplitN(u.Path, "/", -1)
	embed := []*discordgo.MessageEmbed{}

	local := cmd.msgResponse(i.Locale.String())

	if len(arg) != 0 {
		if arg[1] == "tournament" && arg[3] == "event" {
			cmd.slug = arg[1] + "/" + arg[2] + "/" + arg[3] + "/" + arg[4]
			embed = append(embed, cmd.messageEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
				{Name: "**Slug**", Value: cmd.slug},
			}))
		}
	} else {
		embed = append(embed, cmd.messageEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: "**Slug**", Value: local.errorMsg.Input},
		}))
	}

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("setEvent: %v", local.errorMsg.Respond))
	}
}

func (cmd *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options

	format := int(args[0].IntValue())
	stage := args[1].StringValue()
	rounds := int(args[2].IntValue())
	duration := int(args[3].IntValue())

	embed := []*discordgo.MessageEmbed{}

	cmd.rulesMatches.Format = format
	cmd.rulesMatches.Stage = stage
	cmd.rulesMatches.Rounds = rounds
	cmd.rulesMatches.Duration = duration
	cmd.rulesMatches.Crossplatform = args[4].BoolValue()

	local := cmd.msgResponse(i.Locale.String())

	embed = append(embed, cmd.messageEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
		{Name: local.vdMsg.MessageRulesHeader, Value: ""},
		{Name: local.invMsg.Format, Value: fmt.Sprintf(local.invMsg.FT, cmd.tournament.Rules.Format) + fmt.Sprintf(local.invMsg.FormatDescription, cmd.tournament.Rules.Format), Inline: true},
		{Name: local.invMsg.Stage, Value: stage, Inline: true},
		{Name: local.invMsg.Rounds, Value: fmt.Sprintf("%v", cmd.tournament.Rules.Rounds), Inline: true},
		{Name: local.invMsg.Duration, Value: fmt.Sprintf(local.invMsg.DurationCount, cmd.tournament.Rules.Duration), Inline: true},
		{Name: local.invMsg.Crossplatform, Value: local.crossplayRules, Inline: true},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editRuleMatches: %v", local.errorMsg.Respond))
	}
}

func (cmd *commandHandler) editStreamLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options
	area := args[0].StringValue()
	lang := args[1].StringValue()
	conn := args[2].StringValue()
	crossplatform := args[3].BoolValue()
	passcode := args[4].StringValue()
	pc := regexp.MustCompile(`[0-9]+`).FindAllString(passcode, -1)[0]

	embed := []*discordgo.MessageEmbed{}

	local := cmd.msgResponse(i.Locale.String())

	if len(pc) != 4 {
		embed = append(embed, cmd.messageEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: local.vdMsg.MessageStreamHeader, Value: local.errorMsg.Input},
		}))
	} else {
		cmd.streamLobby.Area = area
		cmd.streamLobby.Language = lang
		cmd.streamLobby.Conn = conn
		cmd.streamLobby.Crossplatform = crossplatform
		cmd.streamLobby.Passcode = pc
		embed = append(embed, cmd.messageEmbed(local.vdMsg.MessageStreamHeader, []*discordgo.MessageEmbedField{
			{Name: local.streamMsg.Area, Value: fmt.Sprintf("%v", local.area)},
			{Name: local.streamMsg.Language, Value: local.lang},
			{Name: local.streamMsg.TypeConnection, Value: fmt.Sprintf("%v", local.conn)},
			{Name: local.streamMsg.Crossplatform, Value: fmt.Sprintf("%v", local.crossplayRules)},
			{Name: local.streamMsg.Passcode, Value: fmt.Sprintf(local.streamMsg.PasscodeTemplate, cmd.streamLobby.Passcode)},
		}))
	}

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("editStreamLobby: can't respond on message"))
	}
}

func (cmd *commandHandler) editLogoTournament(s *discordgo.Session, i *discordgo.InteractionCreate) {
	arg := i.ApplicationCommandData().Options[0].StringValue()
	cmd.logoTournament = arg

	embed := []*discordgo.MessageEmbed{}

	local := cmd.msgResponse(i.Locale.String())

	embed = append(embed, cmd.messageEmbed(local.vdMsg.LogoTournament, []*discordgo.MessageEmbedField{
		{Name: "**Url**", Value: fmt.Sprintf("%v", cmd.logoTournament)},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editLogoTournament: %v", local.errorMsg.Respond))
	}
}

func (cmd *commandHandler) getDiscordContacts(s *discordgo.Session) {
	sliceMessages := []*discordgo.MessageEmbed{}
	fields := []*discordgo.MessageEmbedField{}
	counter := 0
	for nickname, dc := range cmd.discordContacts {
		if counter < 25 {
			usr, err := cmd.searchContactDiscord(s, nickname)
			if err != nil {
				log.Printf("viewContacts: %v", err.Error())
				fields = append(fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("%v", nickname), Value: fmt.Sprintf("__Discord:__\n```%v```__GameID:__\n```%v```", dc.discord, dc.gameID), Inline: false,
				})
			} else {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("%v", nickname), Value: fmt.Sprintf("__Discord:__\n<@%v>__GameID:__\n```%v```", usr.discordID, dc.gameID), Inline: false,
				})
			}
			counter++
		} else {
			embed := cmd.messageEmbed("", fields)
			sliceMessages = append(sliceMessages, embed)
			fields = []*discordgo.MessageEmbedField{}
			counter = 0
		}
	}

	embed := cmd.messageEmbed("", fields)
	sliceMessages = append(sliceMessages, embed)
	cmd.embedDiscordContacts = sliceMessages
}

func (cmd *commandHandler) viewContacts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := cmd.msgResponse(i.Locale.String())
	size := len(cmd.discordContacts)
	if len(cmd.embedDiscordContacts) == 0 {
		cts, err := os.ReadFile("contacts.json")
		if err != nil {
			if size != 0 {
				go func() {
					response(s, i, "Loading contacts from csv file...")
				}()

				cmd.getDiscordContacts(s)

				file, err := json.MarshalIndent(cmd.embedDiscordContacts, "", " ")
				if err != nil {
					log.Println(err.Error())
				}

				err = os.WriteFile("contacts.json", file, 0644)
				if err != nil {
					log.Println(err.Error())
				}

				for _, embed := range cmd.embedDiscordContacts {
					if _, err := s.ChannelMessageSendEmbed(i.ChannelID, embed); err != nil {
						log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
					}
				}
			} else {
				embed := []*discordgo.MessageEmbed{}

				embed = append(embed, cmd.messageEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
					{Name: "", Value: local.errorMsg.NoData},
				}))

				if err := cmd.responseEmbed(s, i, embed); err != nil {
					log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
				}
			}
		} else {
			json.Unmarshal(cts, &cmd.embedDiscordContacts)

			for _, embed := range cmd.embedDiscordContacts {
				if _, err := s.ChannelMessageSendEmbed(i.ChannelID, embed); err != nil {
					log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
				}
			}
		}
	} else {
		for _, embed := range cmd.embedDiscordContacts {
			if _, err := s.ChannelMessageSendEmbed(i.ChannelID, embed); err != nil {
				log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
			}
		}
	}
}
