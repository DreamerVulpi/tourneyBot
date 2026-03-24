package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/internal/entity/locale"
)

const (
	ColorDefault = 0x3498db // Blue: standard matches and information messages (Hex: #3498db)
	ColorStream  = 0x9b59b6 // Purple: matches on stream (Hex: #9b59b6)
	ColorFinal   = 0xe74c3c // Red: final stage tournament (Hex: #e74c3c)
	ColorSuccess = 0x2ecc71 // Green: success status (Hex: #2ecc71)
	ColorError   = 0x960018 // Dark-red: errors (Hex: #960018)
	ColorSystem  = 0x9c9c9c // Grey: check data and log (Hex: #9c9c9c)
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

func (_ *discordHandler) typeLocale(language string) locale.Lang {
	var local locale.Lang
	switch language {
	case "Russian":
		local = locale.Ru
	default:
		local = locale.En
	}
	return local
}

func (s *discordHandler) fieldCrossplay(local locale.Lang) string {
	crossplay := local.InviteMessage.CrossplatformStatusTrue
	if !s.cfg.rulesMatches.Crossplatform {
		crossplay = local.InviteMessage.CrossplatformStatusFalse
	}
	return crossplay
}

func (s *discordHandler) fieldStage(local locale.Lang) string {
	stage := local.InviteMessage.AnyStage
	if s.cfg.rulesMatches.Stage != "any" {
		stage = s.cfg.rulesMatches.Stage
	}
	return stage
}
func (s *discordHandler) fieldLanguage(local locale.Lang) string {
	lang := local.StreamLobbyMessage.AnyLanguage
	if s.cfg.streamLobby.Language != "any" {
		lang = local.StreamLobbyMessage.SameLanguage
	}
	return lang
}
func (s *discordHandler) fieldArea(local locale.Lang) string {
	area := local.StreamLobbyMessage.AnyArea
	if s.cfg.streamLobby.Area != "any" {
		area = local.StreamLobbyMessage.CloseArea
	}
	return area
}
func (s *discordHandler) fieldConnection(local locale.Lang) string {
	conn := local.StreamLobbyMessage.AnyConnection
	if s.cfg.streamLobby.Conn != "any" {
		conn = s.cfg.streamLobby.Conn
	}
	return conn
}

func (s *DiscordSender) fieldCrossplay(local locale.Lang) string {
	crossplay := local.InviteMessage.CrossplatformStatusTrue
	if !s.cfg.rulesMatches.Crossplatform {
		crossplay = local.InviteMessage.CrossplatformStatusFalse
	}
	return crossplay
}

func (s *DiscordSender) fieldStage(local locale.Lang) string {
	stage := local.InviteMessage.AnyStage
	if s.cfg.rulesMatches.Stage != "any" {
		stage = s.cfg.rulesMatches.Stage
	}
	return stage
}
func (s *DiscordSender) fieldLanguage(local locale.Lang) string {
	lang := local.StreamLobbyMessage.AnyLanguage
	if s.cfg.streamLobby.Language != "any" {
		lang = local.StreamLobbyMessage.SameLanguage
	}
	return lang
}
func (s *DiscordSender) fieldArea(local locale.Lang) string {
	area := local.StreamLobbyMessage.AnyArea
	if s.cfg.streamLobby.Area != "any" {
		area = local.StreamLobbyMessage.CloseArea
	}
	return area
}
func (s *DiscordSender) fieldConnection(local locale.Lang) string {
	conn := local.StreamLobbyMessage.AnyConnection
	if s.cfg.streamLobby.Conn != "any" {
		conn = s.cfg.streamLobby.Conn
	}
	return conn
}

func (s *DiscordSender) templateEmbedMsg(title string, fields []*discordgo.MessageEmbedField, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: title,
		Color: color,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: s.cfg.logo,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.cfg.tournament.Logo.Img,
		},
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "by DreamerVulpi | https://www.twitch.tv/dreamervulpi",
			IconURL: "https://i.imgur.com/eVmmYEV.png",
		},
	}
}

func (dh *discordHandler) msgStreamLobby(language string, embedColor int) *discordgo.MessageEmbed {
	local := dh.typeLocale(language)

	fields := []*discordgo.MessageEmbedField{
		{Name: local.ViewDataMessage.MessageStreamHeader},
		{Name: local.StreamLobbyMessage.Area, Value: dh.fieldArea(local), Inline: true},
		{Name: local.StreamLobbyMessage.Language, Value: dh.fieldLanguage(local), Inline: true},
		{Name: local.StreamLobbyMessage.TypeConnection, Value: dh.fieldConnection(local), Inline: true},
		{Name: local.StreamLobbyMessage.Crossplatform, Value: dh.fieldCrossplay(local), Inline: true},
		{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, dh.cfg.streamLobby.Passcode), Inline: true},
	}
	message := dh.msgEmbed(local.ViewDataMessage.Title, fields, embedColor)
	return message
}

func (dh *discordHandler) msgRuleMatches(language string, embedColor int) *discordgo.MessageEmbed {
	local := dh.typeLocale(language)

	fields := []*discordgo.MessageEmbedField{
		{Name: local.ViewDataMessage.MessageRulesHeader},
		{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, dh.cfg.rulesMatches.StandardFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, dh.cfg.rulesMatches.StandardFormat), Inline: true},
		{Name: local.InviteMessage.FinalsFormat, Value: fmt.Sprintf(local.InviteMessage.FT, dh.cfg.rulesMatches.FinalsFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, dh.cfg.rulesMatches.FinalsFormat), Inline: true},
		{Name: local.InviteMessage.Stage, Value: dh.fieldStage(local), Inline: true},
		{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", dh.cfg.rulesMatches.Rounds), Inline: true},
		{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, dh.cfg.rulesMatches.Duration), Inline: true},
		{Name: local.InviteMessage.Crossplatform, Value: dh.fieldCrossplay(local), Inline: true},
	}
	message := dh.msgEmbed(local.ViewDataMessage.Title, fields, embedColor)
	return message
}

func (dh *discordHandler) msgViewData(language string) *discordgo.MessageEmbed {
	local := dh.typeLocale(language)

	slug := dh.slug
	if len(slug) == 0 {
		slug = local.ErrorMessage.NoData
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "**Slug**", Value: fmt.Sprintln(local.ViewDataMessage.Description), Inline: true},
		{Value: fmt.Sprintf("```%v```", slug)},

		{Name: local.ViewDataMessage.MessageRulesHeader},
		{Name: local.InviteMessage.StandardFormat,
			Value:  fmt.Sprintf(local.InviteMessage.FT, dh.cfg.rulesMatches.StandardFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, dh.cfg.rulesMatches.StandardFormat),
			Inline: true},
		{Name: local.InviteMessage.FinalsFormat,
			Value:  fmt.Sprintf(local.InviteMessage.FT, dh.cfg.rulesMatches.FinalsFormat) + fmt.Sprintf(local.InviteMessage.FormatDescription, dh.cfg.rulesMatches.FinalsFormat),
			Inline: true},
		{Name: local.InviteMessage.Stage,
			Value:  dh.fieldStage(local),
			Inline: true},
		{Name: local.InviteMessage.Rounds,
			Value:  fmt.Sprintf("%v", dh.cfg.rulesMatches.Rounds),
			Inline: true},
		{Name: local.InviteMessage.Duration,
			Value:  fmt.Sprintf(local.InviteMessage.DurationCount, dh.cfg.rulesMatches.Duration),
			Inline: true},
		{Name: local.InviteMessage.Crossplatform,
			Value:  dh.fieldCrossplay(local),
			Inline: true},

		{Name: local.ViewDataMessage.MessageStreamHeader},
		{Name: local.StreamLobbyMessage.Area,
			Value:  dh.fieldArea(local),
			Inline: true},
		{Name: local.StreamLobbyMessage.Language,
			Value:  dh.fieldLanguage(local),
			Inline: true},
		{Name: local.StreamLobbyMessage.TypeConnection,
			Value:  dh.fieldConnection(local),
			Inline: true},
		{Name: local.StreamLobbyMessage.Crossplatform,
			Value:  dh.fieldCrossplay(local),
			Inline: true},
		{Name: local.StreamLobbyMessage.Passcode,
			Value:  fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, dh.cfg.streamLobby.Passcode),
			Inline: true},
	}
	message := dh.msgEmbed(local.ViewDataMessage.Title, fields, ColorSystem)
	return message
}

func (s *discordHandler) msgResponse(language string) responseLocale {
	local := s.typeLocale(language)

	var result responseLocale
	result.errorMsg = local.ErrorMessage
	result.vdMsg = local.ViewDataMessage
	result.invMsg = local.InviteMessage
	result.streamMsg = local.StreamLobbyMessage
	result.responseMsg = local.ResponseMessage

	rulesCrossplatform := local.InviteMessage.CrossplatformStatusTrue
	if !s.cfg.rulesMatches.Crossplatform {
		rulesCrossplatform = local.InviteMessage.CrossplatformStatusFalse
	}

	streamCrossplatform := local.StreamLobbyMessage.CrossplatformStatusTrue
	if !s.cfg.streamLobby.Crossplatform {
		streamCrossplatform = local.StreamLobbyMessage.CrossplatformStatusFalse
	}

	result.area = s.fieldArea(local)
	result.conn = s.fieldConnection(local)
	result.lang = s.fieldLanguage(local)
	result.crossplayLobby = streamCrossplatform
	result.crossplayRules = rulesCrossplatform

	return result
}

func (s *discordHandler) msgEmbed(title string, fields []*discordgo.MessageEmbedField, color int) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: color,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: s.cfg.logo,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.cfg.tournament.Logo.Img,
		},
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "by DreamerVulpi | https://www.twitch.tv/dreamervulpi",
			IconURL: "https://i.imgur.com/eVmmYEV.png",
		},
	}

	return embed
}
