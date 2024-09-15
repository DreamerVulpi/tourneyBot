package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/locale"
)

// TODO: REFACTOR CODE

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

func (c *commandHandler) msgInviteRu(s *discordgo.Session, player PlayerData, channel *discordgo.Channel, link string) {
	if len(player.streamSourse) == 0 {
		var crossplay string
		if !c.tournament.Rules.Crossplatform {
			crossplay = locale.Ru.InviteMessage.CrossplatformStatusFalse
		} else {
			crossplay = locale.Ru.InviteMessage.CrossplatformStatusTrue
		}
		var stage string
		if c.tournament.Rules.Stage == "any" {
			stage = locale.Ru.InviteMessage.AnyStage
		} else {
			stage = c.tournament.Rules.Stage
		}
		fields := []*discordgo.MessageEmbedField{
			{Name: locale.Ru.InviteMessage.MessageHeader, Value: ""},
			{Name: locale.Ru.InviteMessage.Nickname, Value: fmt.Sprintf("```%v```", player.opponent.nickname), Inline: true},
			{Name: locale.Ru.InviteMessage.TekkenID, Value: fmt.Sprintf("```%v```", player.opponent.tekkenID), Inline: true},
			{Name: locale.Ru.InviteMessage.Discord, Value: fmt.Sprintf("<@%v>", player.opponent.discordID), Inline: true},

			{Name: locale.Ru.InviteMessage.CheckIn, Value: link},
			{Name: fmt.Sprintf(locale.Ru.InviteMessage.Warning, c.tournament.Rules.Waiting), Value: ""},

			{Name: locale.Ru.InviteMessage.SettingsHeader, Value: ""},
			{Name: locale.Ru.InviteMessage.Format, Value: fmt.Sprintf(locale.Ru.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(locale.Ru.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
			{Name: locale.Ru.InviteMessage.Stage, Value: stage, Inline: true},
			{Name: locale.Ru.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
			{Name: locale.Ru.InviteMessage.Duration, Value: fmt.Sprintf(locale.Ru.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
			{Name: locale.Ru.InviteMessage.Crossplatform, Value: crossplay, Inline: true},
		}
		message := c.templateMessage(fields)
		message.Title = fmt.Sprintf(locale.Ru.InviteMessage.Title, player.tournament)
		message.Description = locale.Ru.InviteMessage.Description
		_, err := s.ChannelMessageSendEmbed(channel.ID, message)
		if err != nil {
			log.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				c.m.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		} else {
			var lang string
			if c.tournament.Stream.Language == "any" {
				lang = locale.Ru.StreamLobbyMessage.AnyLanguage
			} else {
				lang = locale.Ru.StreamLobbyMessage.SameLanguage
			}
			var area string
			if c.tournament.Stream.Area == "any" {
				area = locale.Ru.StreamLobbyMessage.AnyArea
			} else {
				area = locale.Ru.StreamLobbyMessage.CloseArea
			}
			var crossplay string
			if !c.tournament.Rules.Crossplatform {
				crossplay = locale.Ru.StreamLobbyMessage.CrossplatformStatusFalse
			} else {
				crossplay = locale.Ru.StreamLobbyMessage.CrossplatformStatusTrue
			}
			var conn string
			if c.tournament.Stream.Conn == "any" {
				conn = locale.Ru.StreamLobbyMessage.AnyConnection
			} else {
				conn = c.tournament.Stream.Conn
			}
			fields := []*discordgo.MessageEmbedField{
				{Name: locale.Ru.StreamLobbyMessage.MessageHeader, Value: ""},
				{Name: fmt.Sprintf(locale.Ru.StreamLobbyMessage.Warning, c.tournament.Rules.Waiting)},
				{Name: "", Value: link, Inline: true},

				{Name: locale.Ru.StreamLobbyMessage.ParamsHeader, Value: ""},
				{Name: "", Value: ""},
				{Name: locale.Ru.StreamLobbyMessage.Area, Value: area, Inline: true},
				{Name: locale.Ru.StreamLobbyMessage.Language, Value: lang, Inline: true},
				{Name: locale.Ru.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
				{Name: locale.Ru.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
				{Name: locale.Ru.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(locale.Ru.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
			}
			message := c.templateMessage(fields)
			message.Title = fmt.Sprintf(locale.Ru.StreamLobbyMessage.Title, player.tournament)
			message.Description = locale.Ru.StreamLobbyMessage.Description
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
}

func (c *commandHandler) msgInviteDefault(s *discordgo.Session, player PlayerData, channel *discordgo.Channel, link string) {
	if len(player.streamSourse) == 0 {
		var crossplay string
		if !c.tournament.Rules.Crossplatform {
			crossplay = locale.En.InviteMessage.CrossplatformStatusFalse
		} else {
			crossplay = locale.En.InviteMessage.CrossplatformStatusTrue
		}
		var stage string
		if c.tournament.Rules.Stage == "any" {
			stage = locale.En.InviteMessage.AnyStage
		} else {
			stage = c.tournament.Rules.Stage
		}
		fields := []*discordgo.MessageEmbedField{
			{Name: locale.En.InviteMessage.MessageHeader, Value: ""},
			{Name: locale.En.InviteMessage.Nickname, Value: fmt.Sprintf("```%v```", player.opponent.nickname), Inline: true},
			{Name: locale.En.InviteMessage.TekkenID, Value: fmt.Sprintf("```%v```", player.opponent.tekkenID), Inline: true},
			{Name: locale.En.InviteMessage.Discord, Value: fmt.Sprintf("<@%v>", player.opponent.discordID), Inline: true},

			{Name: locale.En.InviteMessage.CheckIn, Value: link},
			{Name: fmt.Sprintf(locale.En.InviteMessage.Warning, c.tournament.Rules.Waiting), Value: ""},

			{Name: locale.En.InviteMessage.SettingsHeader, Value: ""},
			{Name: locale.En.InviteMessage.Format, Value: fmt.Sprintf(locale.En.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(locale.En.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
			{Name: locale.En.InviteMessage.Stage, Value: stage, Inline: true},
			{Name: locale.En.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
			{Name: locale.En.InviteMessage.Duration, Value: fmt.Sprintf(locale.En.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
			{Name: locale.En.InviteMessage.Crossplatform, Value: crossplay, Inline: true},
		}
		message := c.templateMessage(fields)
		message.Title = fmt.Sprintf(locale.En.InviteMessage.Title, player.tournament)
		message.Description = locale.En.InviteMessage.Description
		_, err := s.ChannelMessageSendEmbed(channel.ID, message)
		if err != nil {
			log.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				c.m.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		} else {
			var lang string
			if c.tournament.Stream.Language == "any" {
				lang = locale.En.StreamLobbyMessage.AnyLanguage
			} else {
				lang = locale.En.StreamLobbyMessage.SameLanguage
			}
			var area string
			if c.tournament.Stream.Area == "any" {
				area = locale.En.StreamLobbyMessage.AnyArea
			} else {
				area = locale.En.StreamLobbyMessage.CloseArea
			}
			var crossplay string
			if !c.tournament.Rules.Crossplatform {
				crossplay = locale.En.StreamLobbyMessage.CrossplatformStatusFalse
			} else {
				crossplay = locale.En.StreamLobbyMessage.CrossplatformStatusTrue
			}
			var conn string
			if c.tournament.Stream.Conn == "any" {
				conn = locale.En.StreamLobbyMessage.AnyConnection
			} else {
				conn = c.tournament.Stream.Conn
			}
			fields := []*discordgo.MessageEmbedField{
				{Name: locale.En.StreamLobbyMessage.MessageHeader, Value: ""},
				{Name: fmt.Sprintf(locale.En.StreamLobbyMessage.Warning, c.tournament.Rules.Waiting)},
				{Name: "", Value: link, Inline: true},

				{Name: locale.En.StreamLobbyMessage.ParamsHeader, Value: ""},
				{Name: "", Value: ""},
				{Name: locale.En.StreamLobbyMessage.Area, Value: area, Inline: true},
				{Name: locale.En.StreamLobbyMessage.Language, Value: lang, Inline: true},
				{Name: locale.En.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
				{Name: locale.En.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
				{Name: locale.En.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(locale.En.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
			}
			message := c.templateMessage(fields)
			message.Title = fmt.Sprintf(locale.En.StreamLobbyMessage.Title, player.tournament)
			message.Description = locale.En.StreamLobbyMessage.Description
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
}

func (c *commandHandler) msgViewDataDefault() *discordgo.MessageEmbed {
	var crossplay string
	if !c.tournament.Rules.Crossplatform {
		crossplay = locale.En.InviteMessage.CrossplatformStatusFalse
	} else {
		crossplay = locale.En.InviteMessage.CrossplatformStatusTrue
	}
	var stage string
	if c.tournament.Rules.Stage == "any" {
		stage = locale.En.InviteMessage.AnyStage
	} else {
		stage = c.tournament.Rules.Stage
	}
	var lang string
	if c.tournament.Stream.Language == "any" {
		lang = locale.En.StreamLobbyMessage.AnyLanguage
	} else {
		lang = locale.En.StreamLobbyMessage.SameLanguage
	}
	var area string
	if c.tournament.Stream.Area == "any" {
		area = locale.En.StreamLobbyMessage.AnyArea
	} else {
		area = locale.En.StreamLobbyMessage.CloseArea
	}
	var conn string
	if c.tournament.Stream.Conn == "any" {
		conn = locale.En.StreamLobbyMessage.AnyConnection
	} else {
		conn = c.tournament.Stream.Conn
	}
	fields := []*discordgo.MessageEmbedField{
		{Name: locale.En.ViewDataMessage.Title},
		{Name: "**Slug**", Value: fmt.Sprintln(locale.En.ViewDataMessage.Description), Inline: true},
		{Value: fmt.Sprintf("```%v```", c.slug)},

		{Name: locale.En.ViewDataMessage.MessageRulesHeader},
		{Name: locale.En.InviteMessage.Format, Value: fmt.Sprintf(locale.En.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(locale.En.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
		{Name: locale.En.InviteMessage.Stage, Value: stage, Inline: true},
		{Name: locale.En.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
		{Name: locale.En.InviteMessage.Duration, Value: fmt.Sprintf(locale.En.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
		{Name: locale.En.InviteMessage.Crossplatform, Value: crossplay, Inline: true},

		{Name: locale.En.ViewDataMessage.MessageStreamHeader},
		{Name: locale.En.StreamLobbyMessage.Area, Value: area, Inline: true},
		{Name: locale.En.StreamLobbyMessage.Language, Value: lang, Inline: true},
		{Name: locale.En.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
		{Name: locale.En.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
		{Name: locale.En.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(locale.En.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
	}
	message := c.templateMessage(fields)
	return message
}

func (c *commandHandler) msgViewDataRu() *discordgo.MessageEmbed {
	var crossplay string
	if !c.tournament.Rules.Crossplatform {
		crossplay = locale.Ru.InviteMessage.CrossplatformStatusFalse
	} else {
		crossplay = locale.Ru.InviteMessage.CrossplatformStatusTrue
	}
	var stage string
	if c.tournament.Rules.Stage == "any" {
		stage = locale.Ru.InviteMessage.AnyStage
	} else {
		stage = c.tournament.Rules.Stage
	}
	var lang string
	if c.tournament.Stream.Language == "any" {
		lang = locale.Ru.StreamLobbyMessage.AnyLanguage
	} else {
		lang = locale.Ru.StreamLobbyMessage.SameLanguage
	}
	var area string
	if c.tournament.Stream.Area == "any" {
		area = locale.Ru.StreamLobbyMessage.AnyArea
	} else {
		area = locale.Ru.StreamLobbyMessage.CloseArea
	}
	var conn string
	if c.tournament.Stream.Conn == "any" {
		conn = locale.Ru.StreamLobbyMessage.AnyConnection
	} else {
		conn = c.tournament.Stream.Conn
	}
	fields := []*discordgo.MessageEmbedField{
		{Name: locale.Ru.ViewDataMessage.Title},
		{Name: "**Slug**", Value: fmt.Sprintln(locale.Ru.ViewDataMessage.Description), Inline: true},
		{Value: fmt.Sprintf("```%v```", c.slug)},

		{Name: locale.Ru.ViewDataMessage.MessageRulesHeader},
		{Name: locale.Ru.InviteMessage.Format, Value: fmt.Sprintf(locale.Ru.InviteMessage.FT, c.tournament.Rules.Format) + fmt.Sprintf(locale.Ru.InviteMessage.FormatDescription, c.tournament.Rules.Format), Inline: true},
		{Name: locale.Ru.InviteMessage.Stage, Value: stage, Inline: true},
		{Name: locale.Ru.InviteMessage.Rounds, Value: fmt.Sprintf("%v", c.tournament.Rules.Rounds), Inline: true},
		{Name: locale.Ru.InviteMessage.Duration, Value: fmt.Sprintf(locale.Ru.InviteMessage.DurationCount, c.tournament.Rules.Duration), Inline: true},
		{Name: locale.Ru.InviteMessage.Crossplatform, Value: crossplay, Inline: true},

		{Name: locale.Ru.ViewDataMessage.MessageStreamHeader},
		{Name: locale.Ru.StreamLobbyMessage.Area, Value: area, Inline: true},
		{Name: locale.Ru.StreamLobbyMessage.Language, Value: lang, Inline: true},
		{Name: locale.Ru.StreamLobbyMessage.TypeConnection, Value: conn, Inline: true},
		{Name: locale.Ru.StreamLobbyMessage.Crossplatform, Value: crossplay, Inline: true},
		{Name: locale.Ru.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(locale.Ru.StreamLobbyMessage.PasscodeTemplate, c.tournament.Stream.Passcode), Inline: true},
	}
	message := c.templateMessage(fields)
	return message
}

func (c *commandHandler) msgResponse(i *discordgo.InteractionCreate) responseLocale {
	var rslt responseLocale

	if i.Locale.String() == "Russian" {
		rslt.errorMsg = locale.Ru.ErrorMessage
		rslt.vdMsg = locale.Ru.ViewDataMessage
		rslt.invMsg = locale.Ru.InviteMessage
		rslt.responseMsg = locale.Ru.ResponseMessage
		if c.tournament.Rules.Crossplatform {
			rslt.crossplayRules = locale.Ru.InviteMessage.CrossplatformStatusTrue
		} else {
			rslt.crossplayRules = locale.Ru.InviteMessage.CrossplatformStatusFalse
		}
		if c.tournament.Stream.Crossplatform {
			rslt.crossplayLobby = locale.Ru.StreamLobbyMessage.CrossplatformStatusTrue
		} else {
			rslt.crossplayLobby = locale.Ru.StreamLobbyMessage.CrossplatformStatusFalse
		}
		if c.tournament.Stream.Area == "any" {
			rslt.area = locale.Ru.StreamLobbyMessage.AnyArea
		} else {
			rslt.area = locale.Ru.StreamLobbyMessage.CloseArea
		}
		if c.tournament.Stream.Conn == "any" {
			rslt.area = locale.Ru.StreamLobbyMessage.AnyConnection
		} else {
			rslt.area = c.tournament.Stream.Conn
		}
		if c.tournament.Stream.Language == "any" {
			rslt.lang = locale.Ru.StreamLobbyMessage.AnyLanguage
		} else {
			rslt.lang = locale.Ru.StreamLobbyMessage.SameLanguage
		}
	} else {
		rslt.errorMsg = locale.En.ErrorMessage
		rslt.vdMsg = locale.En.ViewDataMessage
		rslt.invMsg = locale.En.InviteMessage

		if c.tournament.Rules.Crossplatform {
			rslt.crossplayRules = locale.En.InviteMessage.CrossplatformStatusTrue
		} else {
			rslt.crossplayRules = locale.En.InviteMessage.CrossplatformStatusFalse
		}

		if c.tournament.Stream.Crossplatform {
			rslt.crossplayLobby = locale.En.StreamLobbyMessage.CrossplatformStatusTrue
		} else {
			rslt.crossplayLobby = locale.En.StreamLobbyMessage.CrossplatformStatusFalse
		}

		if c.tournament.Stream.Area == "any" {
			rslt.area = locale.En.StreamLobbyMessage.AnyArea
		} else {
			rslt.area = locale.En.StreamLobbyMessage.CloseArea
		}
		if c.tournament.Stream.Conn == "any" {
			rslt.area = locale.En.StreamLobbyMessage.AnyConnection
		} else {
			rslt.area = c.tournament.Stream.Conn
		}
		if c.tournament.Stream.Language == "any" {
			rslt.lang = locale.En.StreamLobbyMessage.AnyLanguage
		} else {
			rslt.lang = locale.En.StreamLobbyMessage.SameLanguage
		}
	}
	return rslt
}
