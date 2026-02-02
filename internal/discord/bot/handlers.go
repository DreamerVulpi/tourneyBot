package bot

import (
	"log"
	"sync"

	// "time"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
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
	startgg    strtgg
	discord    discord
	cfg        params
	cancelFunc context.CancelFunc
	mu         sync.Mutex
	debugMode  bool
	auth       *auth.AuthClient
}

func (ch *commandHandler) viewData(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed := []*discordgo.MessageEmbed{
		ch.msgViewData(i.Locale.String()),
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("viewData:", local.errorMsg.Respond)
	}
}

func (ch *commandHandler) startSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())
	if err := ch.processSending(s, i, local); err != nil {
		log.Println("processSending:", err)
	}
}

func (ch *commandHandler) stopSending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	go func() {
		if err := response(s, i, local.responseMsg.Stopping); err != nil {
			log.Println("Error sending interaction response:", err)
		}
	}()

	ch.mu.Lock()
	if ch.cancelFunc != nil {
		ch.cancelFunc()
		ch.cancelFunc = nil
		log.Println("SUCCESS: Cancel function executed")
	}
	ch.mu.Unlock()

	_, err := s.ChannelMessageSend(i.ChannelID, local.responseMsg.Stopped)
	if err != nil {
		log.Println("Error sending final message:", err)
	}
}

func (ch *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())
	embed, err := ch.parseURL(i, local)
	if err != nil {
		log.Println("parseURL:", err)
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("setEvent:", local.errorMsg.Respond)
	}
}

func (ch *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed := ch.getRuleMatchesData(i)
	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("editRuleMatches:", local.errorMsg.Respond)
	}
}

func (ch *commandHandler) editStreamLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed, err := ch.getStreamLobbyData(i, local)
	if err != nil {
		log.Println("getStreamLobbyData:", err)
	}

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("editStreamLobby:", local.errorMsg.Respond)
	}
}

func (ch *commandHandler) editLogoTournament(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	embed := ch.getLogoTournamnentURL(i, local)
	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("editLogoTournament:", local.errorMsg.Respond)
	}
}

func (ch *commandHandler) viewContacts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())

	go func() {
		embed, err := ch.readCommandEmbedJSON(s, i, local)
		if err != nil {
			log.Println("readCommandEmbedJSON:", err)
		}
		if len(embed) > 0 {
			if err := ch.responseEmbed(s, i, embed); err != nil {
				log.Println("viewContacts response error:", err)
			}
		}
	}()
}

func (ch *commandHandler) roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	local := ch.msgResponse(i.Locale.String())
	arg := i.ApplicationCommandData().Options[0].StringValue()

	embed := ch.controlRole(s, arg)

	if err := ch.responseEmbed(s, i, embed); err != nil {
		log.Println("roles:", local.errorMsg.Respond)
	}
}
