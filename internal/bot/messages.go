package bot

import (
	"fmt"
	"log"

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

func (c *commandHandler) msgInvite(s *discordgo.Session, player PlayerData, channel *discordgo.Channel, link string, roleId string) {
	var local locale.Lang
	if roleId == c.rolesIdList.Ru {
		local = locale.Ru
	} else {
		local = locale.En
	}

	if len(player.streamSourse) == 0 {
		crossplay := local.InviteMessage.CrossplatformStatusTrue
		if !c.tournament.Rules.Crossplatform {
			crossplay = local.InviteMessage.CrossplatformStatusFalse
		}
		stage := local.InviteMessage.AnyStage
		if c.tournament.Rules.Stage != "any" {
			stage = c.tournament.Rules.Stage
		}
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
			{Name: fmt.Sprintf(local.InviteMessage.Warning, c.tournament.Rules.Waiting), Value: ""},

			{Name: local.InviteMessage.SettingsHeader},
			{Name: local.InviteMessage.Format, Value: fmt.Sprintf(local.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(local.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
			{Name: local.InviteMessage.Stage, Value: stage, Inline: true},
			{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
			{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
			{Name: local.InviteMessage.Crossplatform, Value: crossplay, Inline: true},
		}
		message := c.templateMessage(fields)
		message.Title = fmt.Sprintf(local.InviteMessage.Title, player.tournament)
		message.Description = local.InviteMessage.Description
		_, err := s.ChannelMessageSendEmbed(channel.ID, message)
		if err != nil {
			log.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				c.m.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		}
	} else {
		lang := local.StreamLobbyMessage.AnyLanguage
		if c.tournament.Stream.Language != "any" {
			lang = local.StreamLobbyMessage.SameLanguage
		}
		area := local.StreamLobbyMessage.AnyArea
		if c.tournament.Stream.Area != "any" {
			area = local.StreamLobbyMessage.CloseArea
		}
		crossplay := local.StreamLobbyMessage.CrossplatformStatusTrue
		if !c.tournament.Rules.Crossplatform {
			crossplay = local.StreamLobbyMessage.CrossplatformStatusFalse
		}
		conn := local.StreamLobbyMessage.AnyConnection
		if c.tournament.Stream.Conn != "any" {
			conn = c.tournament.Stream.Conn
		}

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
			{Name: fmt.Sprintf(local.StreamLobbyMessage.Warning, c.tournament.Rules.Waiting)},

			{Name: local.StreamLobbyMessage.ParamsHeader},
			{Name: local.StreamLobbyMessage.Area, Value: area, Inline: true},
			{Name: local.StreamLobbyMessage.Language, Value: lang, Inline: true},
			{Name: local.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
			{Name: local.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
			{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
		}
		message := c.templateMessage(fields)
		message.Title = fmt.Sprintf(local.StreamLobbyMessage.Title, player.tournament)
		message.Description = local.StreamLobbyMessage.Description
		_, err := s.ChannelMessageSendEmbed(channel.ID, message)
		if err != nil {
			log.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				c.m.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		}
	}
}

func (c *commandHandler) msgViewData(language string) *discordgo.MessageEmbed {
	var local locale.Lang
	switch language {
	case "Russian":
		local = locale.Ru
	default:
		local = locale.En
	}

	crossplay := local.InviteMessage.CrossplatformStatusTrue
	if !c.tournament.Rules.Crossplatform {
		crossplay = local.InviteMessage.CrossplatformStatusFalse
	}
	stage := local.InviteMessage.AnyStage
	if c.tournament.Rules.Stage != "any" {
		stage = c.tournament.Rules.Stage
	}
	lang := local.StreamLobbyMessage.AnyLanguage
	if c.tournament.Stream.Language != "any" {
		lang = local.StreamLobbyMessage.SameLanguage
	}
	area := local.StreamLobbyMessage.AnyArea
	if c.tournament.Stream.Area != "any" {
		area = local.StreamLobbyMessage.CloseArea
	}
	conn := local.StreamLobbyMessage.AnyConnection
	if c.tournament.Stream.Conn != "any" {
		conn = c.tournament.Stream.Conn
	}

	slug := c.slug
	if len(slug) == 0 {
		slug = local.ErrorMessage.NoData
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: local.ViewDataMessage.Title},
		{Name: "**Slug**", Value: fmt.Sprintln(local.ViewDataMessage.Description), Inline: true},
		{Value: fmt.Sprintf("```%v```", slug)},

		{Name: local.ViewDataMessage.MessageRulesHeader},
		{Name: local.InviteMessage.Format, Value: fmt.Sprintf(local.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(local.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
		{Name: local.InviteMessage.Stage, Value: stage, Inline: true},
		{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
		{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
		{Name: local.InviteMessage.Crossplatform, Value: crossplay, Inline: true},

		{Name: local.ViewDataMessage.MessageStreamHeader},
		{Name: local.StreamLobbyMessage.Area, Value: area, Inline: true},
		{Name: local.StreamLobbyMessage.Language, Value: lang, Inline: true},
		{Name: local.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
		{Name: local.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
		{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
	}
	message := c.templateMessage(fields)
	return message
}

func (c *commandHandler) msgResponse(language string) responseLocale {
	var local locale.Lang
	switch language {
	case "Russian":
		local = locale.Ru
	default:
		local = locale.En
	}

	var result responseLocale
	result.errorMsg = local.ErrorMessage
	result.vdMsg = local.ViewDataMessage
	result.invMsg = local.InviteMessage
	result.streamMsg = local.StreamLobbyMessage
	result.responseMsg = local.ResponseMessage

	rulesCrossplatform := local.InviteMessage.CrossplatformStatusTrue
	if !c.tournament.Rules.Crossplatform {
		rulesCrossplatform = local.InviteMessage.CrossplatformStatusFalse
	}

	streamCrossplatform := local.StreamLobbyMessage.CrossplatformStatusTrue
	if !c.tournament.Stream.Crossplatform {
		streamCrossplatform = local.StreamLobbyMessage.CrossplatformStatusFalse
	}

	area := local.StreamLobbyMessage.AnyArea
	if c.tournament.Stream.Area != "any" {
		area = local.StreamLobbyMessage.CloseArea
	}

	conn := local.StreamLobbyMessage.AnyConnection
	if c.tournament.Stream.Conn != "any" {
		conn = c.tournament.Stream.Conn
	}

	lang := local.StreamLobbyMessage.AnyLanguage
	if c.tournament.Stream.Language != "any" {
		lang = local.StreamLobbyMessage.SameLanguage
	}

	result.area = area
	result.conn = conn
	result.lang = lang
	result.crossplayLobby = streamCrossplatform
	result.crossplayRules = rulesCrossplatform

	return result
}
