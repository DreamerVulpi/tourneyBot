package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/locale"
)

type responseLocale struct {
	errorMsg       locale.ErrorMessage
	vdMsg          locale.ViewDataMessage
	invMsg         locale.InviteMessage
	streamMsg      locale.StreamLobbyMessage
	responseMsg    locale.ResponseMessage
	crossplayRules string
	crossplayLobby string
	area           string
	lang           string
	conn           string
}

func (ch *commandHandler) typeLocale(language string) locale.Lang {
	var local locale.Lang
	switch language {
	case "Russian":
		local = locale.Ru
	default:
		local = locale.En
	}
	return local
}

func (ch *commandHandler) fieldCrossplay(local locale.Lang) string {
	crossplay := local.InviteMessage.CrossplatformStatusTrue
	if !ch.cfg.rulesMatches.Crossplatform {
		crossplay = local.InviteMessage.CrossplatformStatusFalse
	}
	return crossplay
}

func (ch *commandHandler) fieldStage(local locale.Lang) string {
	stage := local.InviteMessage.AnyStage
	if ch.cfg.rulesMatches.Stage != "any" {
		stage = ch.cfg.rulesMatches.Stage
	}
	return stage
}
func (ch *commandHandler) fieldLanguage(local locale.Lang) string {
	lang := local.StreamLobbyMessage.AnyLanguage
	if ch.cfg.streamLobby.Language != "any" {
		lang = local.StreamLobbyMessage.SameLanguage
	}
	return lang
}
func (ch *commandHandler) fieldArea(local locale.Lang) string {
	area := local.StreamLobbyMessage.AnyArea
	if ch.cfg.streamLobby.Area != "any" {
		area = local.StreamLobbyMessage.CloseArea
	}
	return area
}
func (ch *commandHandler) fieldConnection(local locale.Lang) string {
	conn := local.StreamLobbyMessage.AnyConnection
	if ch.cfg.streamLobby.Conn != "any" {
		conn = ch.cfg.streamLobby.Conn
	}
	return conn
}

func (ch *commandHandler) msgInvite(s *discordgo.Session, player PlayerData, channel *discordgo.Channel, link string, roleId string) {
	var message *discordgo.MessageEmbed
	var local locale.Lang
	if roleId == ch.cfg.rolesIdList.Ru {
		local = locale.Ru
	} else {
		local = locale.En
	}

	var format int = ch.cfg.rulesMatches.StandardFormat
	if ch.startgg.finalBracketId == player.phaseGroupId {
		if ch.startgg.minRoundNumA <= player.roundNum && player.roundNum <= ch.startgg.minRoundNumB || ch.startgg.maxRoundNumA <= player.roundNum && player.roundNum <= ch.startgg.maxRoundNumB {
			format = ch.cfg.rulesMatches.FinalsFormat
		}
	}

	if len(player.streamSourse) == 0 {
		gameID := player.opponent.gameID
		if len(gameID) == 0 {
			gameID = local.ErrorMessage.NoData
		}
		discordId := fmt.Sprintf("<@%v>", player.opponent.discordID)
		if len(discordId) == 0 {
			discordId = local.ErrorMessage.NoData
		}

		fields := []*discordgo.MessageEmbedField{
			{Name: local.InviteMessage.MessageHeader},
			{Name: local.InviteMessage.Nickname, Value: fmt.Sprintf("```%v```", player.opponent.nickname), Inline: true},
			{Name: local.InviteMessage.GameID, Value: fmt.Sprintf("```%v```", gameID), Inline: true},
			{Name: local.InviteMessage.Discord, Value: discordId, Inline: true},

			{Name: local.InviteMessage.CheckIn, Value: link},
			{Name: fmt.Sprintf(local.InviteMessage.Warning, ch.cfg.rulesMatches.Waiting), Value: ""},

			{Name: local.InviteMessage.SettingsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.InviteMessage.Stage, Value: ch.fieldStage(local), Inline: true},
			{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", ch.cfg.rulesMatches.Rounds), Inline: true},
			{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, ch.cfg.rulesMatches.Duration), Inline: true},
			{Name: local.InviteMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},
		}
		message = ch.msgEmbed(fmt.Sprintf(local.InviteMessage.Title, player.tournament), fields)
		message.Description = local.InviteMessage.Description
	} else {
		var stream string
		if player.streamSourse == "TWITCH" {
			stream = "https://www.twitch.tv/" + player.streamName
		}
		if player.streamSourse == "YOUTUBE" {
			stream = "https://www.youtube.com/@" + player.streamName
		}
		fields := []*discordgo.MessageEmbedField{
			{Name: local.StreamLobbyMessage.StreamLink, Value: stream},
			{Name: local.StreamLobbyMessage.MessageHeader, Value: link},
			{Name: fmt.Sprintf(local.StreamLobbyMessage.Warning, ch.cfg.rulesMatches.Waiting)},

			{Name: local.StreamLobbyMessage.ParamsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.StreamLobbyMessage.Area, Value: ch.fieldArea(local), Inline: true},
			{Name: local.StreamLobbyMessage.Language, Value: ch.fieldLanguage(local), Inline: true},
			{Name: local.StreamLobbyMessage.TypeConnection, Value: ch.fieldConnection(local), Inline: true},
			{Name: local.StreamLobbyMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},
			{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, ch.cfg.streamLobby.Passcode), Inline: true},
		}
		message = ch.msgEmbed(fmt.Sprintf(local.StreamLobbyMessage.Title, player.tournament), fields)
		message.Description = local.StreamLobbyMessage.Description
	}
	_, err := s.ChannelMessageSendEmbed(channel.ID, message)
	if err != nil {
		log.Println("error sending DM message:", err)
		if _, err := s.ChannelMessageSend(
			ch.discord.msgCreate.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		); err != nil {
			log.Println(err.Error())
		}
	}
}

func (ch *commandHandler) msgRuleMatches(language string) *discordgo.MessageEmbed {
	local := ch.typeLocale(language)

	fields := []*discordgo.MessageEmbedField{
		{Name: local.ViewDataMessage.MessageRulesHeader},
		{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, ch.cfg.rulesMatches.StandardFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, ch.cfg.rulesMatches.StandardFormat), Inline: true},
		{Name: local.InviteMessage.FinalsFormat, Value: fmt.Sprintf(local.InviteMessage.FT, ch.cfg.rulesMatches.FinalsFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, ch.cfg.rulesMatches.FinalsFormat), Inline: true},
		{Name: local.InviteMessage.Stage, Value: ch.fieldStage(local), Inline: true},
		{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", ch.cfg.rulesMatches.Rounds), Inline: true},
		{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, ch.cfg.rulesMatches.Duration), Inline: true},
		{Name: local.InviteMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},
	}
	message := ch.msgEmbed(local.ViewDataMessage.Title, fields)
	return message
}

func (ch *commandHandler) msgStreamLobby(language string) *discordgo.MessageEmbed {
	local := ch.typeLocale(language)

	fields := []*discordgo.MessageEmbedField{
		{Name: local.ViewDataMessage.MessageStreamHeader},
		{Name: local.StreamLobbyMessage.Area, Value: ch.fieldArea(local), Inline: true},
		{Name: local.StreamLobbyMessage.Language, Value: ch.fieldLanguage(local), Inline: true},
		{Name: local.StreamLobbyMessage.TypeConnection, Value: ch.fieldConnection(local), Inline: true},
		{Name: local.StreamLobbyMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},
		{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, ch.cfg.streamLobby.Passcode), Inline: true},
	}
	message := ch.msgEmbed(local.ViewDataMessage.Title, fields)
	return message
}

func (ch *commandHandler) msgViewData(language string) *discordgo.MessageEmbed {
	local := ch.typeLocale(language)

	slug := ch.slug
	if len(slug) == 0 {
		slug = local.ErrorMessage.NoData
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "**Slug**", Value: fmt.Sprintln(local.ViewDataMessage.Description), Inline: true},
		{Value: fmt.Sprintf("```%v```", slug)},

		{Name: local.ViewDataMessage.MessageRulesHeader},
		{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, ch.cfg.rulesMatches.StandardFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, ch.cfg.rulesMatches.StandardFormat), Inline: true},
		{Name: local.InviteMessage.FinalsFormat, Value: fmt.Sprintf(local.InviteMessage.FT, ch.cfg.rulesMatches.FinalsFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, ch.cfg.rulesMatches.FinalsFormat), Inline: true},
		{Name: local.InviteMessage.Stage, Value: ch.fieldStage(local), Inline: true},
		{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", ch.cfg.rulesMatches.Rounds), Inline: true},
		{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, ch.cfg.rulesMatches.Duration), Inline: true},
		{Name: local.InviteMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},

		{Name: local.ViewDataMessage.MessageStreamHeader},
		{Name: local.StreamLobbyMessage.Area, Value: ch.fieldArea(local), Inline: true},
		{Name: local.StreamLobbyMessage.Language, Value: ch.fieldLanguage(local), Inline: true},
		{Name: local.StreamLobbyMessage.TypeConnection, Value: ch.fieldConnection(local), Inline: true},
		{Name: local.StreamLobbyMessage.Crossplatform, Value: ch.fieldCrossplay(local), Inline: true},
		{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, ch.cfg.streamLobby.Passcode), Inline: true},
	}
	message := ch.msgEmbed(local.ViewDataMessage.Title, fields)
	return message
}

func (ch *commandHandler) msgResponse(language string) responseLocale {
	local := ch.typeLocale(language)

	var result responseLocale
	result.errorMsg = local.ErrorMessage
	result.vdMsg = local.ViewDataMessage
	result.invMsg = local.InviteMessage
	result.streamMsg = local.StreamLobbyMessage
	result.responseMsg = local.ResponseMessage

	rulesCrossplatform := local.InviteMessage.CrossplatformStatusTrue
	if !ch.cfg.rulesMatches.Crossplatform {
		rulesCrossplatform = local.InviteMessage.CrossplatformStatusFalse
	}

	streamCrossplatform := local.StreamLobbyMessage.CrossplatformStatusTrue
	if !ch.cfg.streamLobby.Crossplatform {
		streamCrossplatform = local.StreamLobbyMessage.CrossplatformStatusFalse
	}

	result.area = ch.fieldArea(local)
	result.conn = ch.fieldConnection(local)
	result.lang = ch.fieldLanguage(local)
	result.crossplayLobby = streamCrossplatform
	result.crossplayRules = rulesCrossplatform

	return result
}

func (ch *commandHandler) msgEmbed(title string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: title,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: ch.cfg.logo,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ch.cfg.tournament.Logo.Img,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "by DreamerVulpi | https://www.twitch.tv/dreamervulpi",
			IconURL: "https://i.imgur.com/FcuAfRw.png",
		},
	}
}
