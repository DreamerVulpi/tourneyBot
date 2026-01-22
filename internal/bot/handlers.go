package bot

import (
	"errors"
	"fmt"
	"log"

	// "time"

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
	embed := []*discordgo.MessageEmbed{
		ch.msgViewData(i.Locale.String()),
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("viewData: can't respond on message"))
	}
}

func (ch *commandHandler) startSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())
	embed := []*discordgo.MessageEmbed{
		ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "", Value: local.errorMsg.Input},
		})}

	if err := ch.processSending(s, i, local); err != nil {
		log.Printf("processSending: %v", err)
	}
	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Printf("responseEmbed: %v\n", local.errorMsg.Respond)
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
	local := ch.msgResponse(i.Locale.String())
	embed, err := ch.parseURL(i, local)
	if err != nil {
		log.Println(fmt.Errorf("parseURL: %v\n", err))
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("setEvent: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed := ch.getRuleMatchesData(i, local)
	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editRuleMatches: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) editStreamLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed, err := ch.getStreamLobbyData(i, local)
	if err != nil {
		log.Printf("getStreamLobbyData: %v\n", err)
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editStreamLobby: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) editLogoTournament(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed := ch.getLogoTournamnentURL(i, local)
	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("editLogoTournament: %v", local.errorMsg.Respond))
	}
}

func (ch *commandHandler) viewContacts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	go func() {
		embed, err := ch.readCommandEmbedJSON(s, i, local)
		if err != nil {
			log.Println(fmt.Errorf("readCommandEmbedJSON: %v\n", err))
		}
		if err := ch.responseEmbed(s, i, embed); err != nil {
			log.Println(fmt.Errorf("viewContacts: %v", local.errorMsg.Respond))
		}
	}()
}

func (ch *commandHandler) roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())
	arg := i.ApplicationCommandData().Options[0].StringValue()

	embed := ch.controlRole(s, arg)

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println(fmt.Errorf("roles: %v", local.errorMsg.Respond))
	}
}
