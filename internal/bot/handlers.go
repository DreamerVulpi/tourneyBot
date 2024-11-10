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

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type discord struct {
	msgCreate     *discordgo.MessageCreate
	contacts      map[string]contactData
	embedContacts []*discordgo.MessageEmbed
	tourneyRole   *discordgo.Role
}

type strtgg struct {
	client         *startgg.Client
	finalBracketId int64
	minRoundNumA   int
	minRoundNumB   int
	maxRoundNumA   int
	maxRoundNumB   int
}

type commandHandler struct {
	slug       string
	stopSignal chan struct{}
	startgg    strtgg
	discord    discord
	cfg        params
}

func (ch *commandHandler) viewData(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, ch.msgViewData(i.Locale.String()))

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("check: can't respond on message"))
	}
}

func (ch *commandHandler) startSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	if ch.cfg.guildID != "" && ch.slug != "" {
		if err := response(s, i, local.responseMsg.Starting); err != nil {
			log.Println(err.Error())
		}
		go ch.Process(s)
	} else {
		embed := []*discordgo.MessageEmbed{}

		embed = append(embed, ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "", Value: local.errorMsg.Input},
		}))

		if err := ch.responseEmbed(s, i, embed); err != nil {
			log.Println(fmt.Errorf("editLogoTournament: %v", local.errorMsg.Respond))
		}
	}
}

func (ch *commandHandler) stopSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	go func() {
		if err := response(s, i, local.responseMsg.Stopping); err != nil {
			log.Println(err.Error())
			if _, err := s.ChannelMessageSend(i.ChannelID, err.Error()); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	// Send signal to stop process
	ch.stopSignal <- struct{}{}

	_, err := s.ChannelMessageSend(i.ChannelID, local.responseMsg.Stopped)
	if err != nil {
		log.Println(err.Error())
	}
}

func (ch *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	u, err := url.Parse(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		log.Println(err)
	}

	arg := strings.SplitN(u.Path, "/", -1)
	embed := []*discordgo.MessageEmbed{}

	local := ch.msgResponse(i.Locale.String())

	if len(arg) != 0 {
		if arg[1] == "tournament" && arg[3] == "event" {
			ch.slug = arg[1] + "/" + arg[2] + "/" + arg[3] + "/" + arg[4]
			embed = append(embed, ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
				{Name: "**Slug**", Value: ch.slug},
			}))
		}
	} else {
		embed = append(embed, ch.msgEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: "**Slug**", Value: local.errorMsg.Input},
		}))
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("setEvent: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options
	ch.cfg.rulesMatches.StandardFormat = int(args[0].IntValue())
	ch.cfg.rulesMatches.FinalsFormat = int(args[1].IntValue())
	ch.cfg.rulesMatches.Stage = args[2].StringValue()
	ch.cfg.rulesMatches.Rounds = int(args[3].IntValue())
	ch.cfg.rulesMatches.Duration = int(args[4].IntValue())
	ch.cfg.rulesMatches.Crossplatform = args[5].BoolValue()

	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, ch.msgRuleMatches(i.Locale.String()))

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("check: can't respond on message"))
	}
}

func (ch *commandHandler) editStreamLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options
	pc := regexp.MustCompile(`[0-9]+`).FindAllString(args[4].StringValue(), -1)[0]

	local := ch.msgResponse(i.Locale.String())
	embed := []*discordgo.MessageEmbed{}

	if len(pc) != 4 {
		embed = append(embed, ch.msgEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: local.vdMsg.MessageStreamHeader, Value: local.errorMsg.Input},
		}))
	} else {
		ch.cfg.streamLobby.Area = args[0].StringValue()
		ch.cfg.streamLobby.Language = args[1].StringValue()
		ch.cfg.streamLobby.Conn = args[2].StringValue()
		ch.cfg.streamLobby.Crossplatform = args[3].BoolValue()
		ch.cfg.streamLobby.Passcode = pc
		embed = append(embed, ch.msgStreamLobby(i.Locale.String()))
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("editStreamLobby: can't respond on message"))
	}
}

func (ch *commandHandler) editLogoTournament(s *discordgo.Session, i *discordgo.InteractionCreate) {
	arg := i.ApplicationCommandData().Options[0].StringValue()
	ch.cfg.tournament.Game.Name = arg

	embed := []*discordgo.MessageEmbed{}

	local := ch.msgResponse(i.Locale.String())

	embed = append(embed, ch.msgEmbed(local.vdMsg.LogoTournament, []*discordgo.MessageEmbedField{
		{Name: "**Url**", Value: fmt.Sprintf("%v", ch.cfg.tournament.Game.Name)},
	}))

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editLogoTournament: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) viewContacts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	errRespond := func(embed []*discordgo.MessageEmbed, typeRespond string) {
		embed = append(embed, ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "", Value: local.errorMsg.NoData},
		}))

		if err := ch.responseEmbed(s, i, embed); err != nil {
			log.Println(fmt.Errorf("viewContacts: %v | %v", typeRespond, err.Error()))
		}
	}

	go func() {
		embed := []*discordgo.MessageEmbed{}

		cts, err := os.ReadFile("contactsEmbed.json")
		if err != nil {
			log.Println(err.Error())
			errRespond(embed, local.errorMsg.Respond)
		} else {
			if err := json.Unmarshal(cts, &ch.discord.embedContacts); err != nil {
				log.Println(err.Error())
				errRespond(embed, local.errorMsg.Respond)
			}

			arg := i.ApplicationCommandData().Options[0].StringValue()
			switch strings.ToLower(arg) {
			case "any", "все":
				if err := response(s, i, local.responseMsg.InProcess); err != nil {
					log.Println(err.Error())
				}
				for _, embedContact := range ch.discord.embedContacts {
					if _, err := s.ChannelMessageSendEmbed(i.ChannelID, embedContact); err != nil {
						log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
					}
				}
			default:
				var trigger bool
				for _, embedContact := range ch.discord.embedContacts {
					for _, field := range embedContact.Fields {
						if strings.EqualFold(arg, field.Name) {
							trigger = true
							var fields []*discordgo.MessageEmbedField
							fields = append(fields, &discordgo.MessageEmbedField{
								Name:  field.Name,
								Value: field.Value,
							})
							embed = append(embed, ch.msgEmbed(local.vdMsg.Title, fields))
							if err := ch.responseEmbed(s, i, embed); err != nil {
								log.Println(fmt.Errorf("viewContacts: %v | %v", local.errorMsg.Respond, err.Error()))
							}
							break
						}
					}
				}
				if !trigger {
					errRespond(embed, local.errorMsg.Input)
				}
			}
		}
	}()
}

func (ch *commandHandler) roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	arg := i.ApplicationCommandData().Options[0].StringValue()

	embed := ch.workRoles(s, arg)

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("roles: %v", local.errorMsg.Respond))
	}
}
