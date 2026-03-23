package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
	dbUC "github.com/dreamervulpi/tourneyBot/internal/db/usecase"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	senderUC "github.com/dreamervulpi/tourneyBot/internal/sender/usecase"
	"github.com/dreamervulpi/tourneyBot/locale"
)

type DiscordSender struct {
	session       *discordgo.Session
	cfg           params
	participantUC dbUC.Participant
	adminID       string
	debugMode     bool
}

func (s *DiscordSender) logSentMsgToDiscord(success bool, errStr string, set sender.SetData, local locale.Lang, gameNickname string) {
	var logFields []*discordgo.MessageEmbedField
	var color int

	var state string
	if success {
		color = ColorSuccess
		state = local.LogMessage.SuccessfulMsgHeader
	} else {
		color = ColorError
		state = local.LogMessage.FailMsgHeader
	}

	logFields = []*discordgo.MessageEmbedField{
		{Name: fmt.Sprintf(local.LogMessage.SuccesfullSendedMsg, state), Value: "\u200B"},
		{Name: fmt.Sprintf("Set #%v | ", set.SetID), Value: fmt.Sprintf("%v vs %v", set.ContactPlayer1.GameNickname, set.ContactPlayer2.GameNickname)},
	}

	if len(set.FullInviteLink) > 0 {
		logFields = append(logFields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf(local.LogMessage.CheckIn, set.FullInviteLink), Value: "\u200B",
		})
	}

	if !success {
		logFields = append(logFields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf(local.LogMessage.FailedSentMsg, gameNickname, errStr), Value: "\u200B",
		})
	}

	logEmbed := s.templateEmbedMsg(fmt.Sprintf(local.LogMessage.Title, set.TournamentName), logFields, color)
	if _, err := s.session.ChannelMessageSendEmbed(s.cfg.debugChannelID, logEmbed); err != nil {
		log.Printf("logToDiscord | error sending to debug channel: %v", err)
	}
}

func (s *DiscordSender) prepareSetData(recipient, opponent sender.Participant, set sender.SetData, local locale.Lang) (*discordgo.MessageEmbed, error) {
	format := s.cfg.rulesMatches.StandardFormat
	embedColor := ColorDefault

	if set.IsFinals {
		format = s.cfg.rulesMatches.FinalsFormat
		embedColor = ColorFinal
	}

	if len(set.StreamSourse) > 0 {
		embedColor = ColorStream
	}

	var message *discordgo.MessageEmbed
	gameNickname := opponent.GameNickname
	if len(gameNickname) == 0 || gameNickname == "N/D" {
		gameNickname = local.ErrorMessage.NoData
	}

	gameID := opponent.GameID
	if len(gameID) == 0 {
		gameID = local.ErrorMessage.NoData
	}

	rawID := opponent.MessenagerID
	login := opponent.MessenagerLogin

	var discordDisplay string
	if len(rawID) == 5 && len(login) == 0 {
		discordDisplay = local.ErrorMessage.NoData
	} else {
		if len(rawID) > 0 && rawID != "000000000000000000" && rawID != "N/D" {
			discordDisplay = fmt.Sprintf("<@%v>", rawID)
		} else {
			if len(login) > 0 {
				discordDisplay = login
			} else {
				discordDisplay = local.ErrorMessage.NoData
			}
		}
	}
	log.Printf("prepareSetData | Set: %v | Recipient: %s vs Opponent: %s | Link: %s", set.SetID, recipient.MessenagerLogin, opponent.MessenagerLogin, set.FullInviteLink)

	if len(set.StreamSourse) == 0 {
		fields := []*discordgo.MessageEmbedField{
			{Name: local.InviteMessage.MessageHeader},
			{Name: local.InviteMessage.Nickname, Value: fmt.Sprintf("```%v```", gameNickname), Inline: true},
			{Name: local.InviteMessage.GameID, Value: fmt.Sprintf("```%v```", gameID), Inline: true},
			{Name: local.InviteMessage.Discord, Value: discordDisplay, Inline: true},

			{Name: local.InviteMessage.CheckIn, Value: set.FullInviteLink},
			{Name: fmt.Sprintf(local.InviteMessage.Warning, s.cfg.rulesMatches.Waiting), Value: "\u200B"},

			{Name: local.InviteMessage.SettingsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.InviteMessage.Stage, Value: s.fieldStage(local), Inline: true},
			{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", s.cfg.rulesMatches.Rounds), Inline: true},
			{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, s.cfg.rulesMatches.Duration), Inline: true},
			{Name: local.InviteMessage.Crossplatform, Value: s.fieldCrossplay(local), Inline: true},
		}
		message = s.templateEmbedMsg(fmt.Sprintf(local.InviteMessage.Title, set.TournamentName), fields, embedColor)
		message.Description = local.InviteMessage.Description
	} else {
		var stream string
		if set.StreamSourse == "TWITCH" {
			stream = "https://www.twitch.tv/" + set.StreamName
		}
		if set.StreamSourse == "YOUTUBE" {
			stream = "https://www.youtube.com/@" + set.StreamName
		}
		fields := []*discordgo.MessageEmbedField{
			{Name: local.StreamLobbyMessage.StreamLink, Value: stream},
			{Name: local.StreamLobbyMessage.MessageHeader, Value: set.FullInviteLink},
			{Name: fmt.Sprintf(local.StreamLobbyMessage.Warning, s.cfg.rulesMatches.Waiting), Value: "\u200B"},

			{Name: local.StreamLobbyMessage.ParamsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.StreamLobbyMessage.Area, Value: s.fieldArea(local), Inline: true},
			{Name: local.StreamLobbyMessage.Language, Value: s.fieldLanguage(local), Inline: true},
			{Name: local.StreamLobbyMessage.TypeConnection, Value: s.fieldConnection(local), Inline: true},
			{Name: local.StreamLobbyMessage.Crossplatform, Value: s.fieldCrossplay(local), Inline: true},
			{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, s.cfg.streamLobby.Passcode), Inline: true},
		}
		message = s.templateEmbedMsg(fmt.Sprintf(local.StreamLobbyMessage.Title, set.TournamentName), fields, embedColor)
		message.Description = local.StreamLobbyMessage.Description
	}
	return message, nil
}

func (s *DiscordSender) msgInvite(targetID string, set sender.SetData, channel *discordgo.Channel) {
	var recipient sender.Participant
	var opponent sender.Participant
	var sidePrefix string

	if targetID == set.ContactPlayer1.MessenagerID {
		recipient = set.ContactPlayer1
		opponent = set.ContactPlayer2
		sidePrefix = "[P1] "
	} else {
		recipient = set.ContactPlayer2
		opponent = set.ContactPlayer1
		sidePrefix = "[P2] "
	}

	// TODO: Change reconize locale in future
	local := locale.En
	if len(recipient.Locales) > 0 {
		local = locale.Ru
	}

	message, err := s.prepareSetData(recipient, opponent, set, local)
	if err != nil {
		log.Printf("msgInvite | error sended DM: %v\n", err.Error())
		s.logSentMsgToDiscord(false, err.Error(), set, local, recipient.GameNickname)
		return
	}

	if set.IsTest {
		message.Title = sidePrefix + message.Title
	}

	_, err = s.session.ChannelMessageSendEmbed(channel.ID, message)
	if err != nil {
		log.Printf("msgInvite | error sended DM: %v\n", err.Error())
		s.logSentMsgToDiscord(false, err.Error(), set, local, recipient.GameNickname)
		return
	}

	s.logSentMsgToDiscord(true, "", set, local, recipient.GameNickname)
	log.Printf("msgInvite | success sended DM to %s (%s)", recipient.GameNickname, sidePrefix)
}

func (s *DiscordSender) SendNotification(ctx context.Context, targetID string, set sender.SetData) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if targetID == "" || targetID == "N/D" || len(targetID) == 0 || targetID == "0" {
		return fmt.Errorf("SendNotification | tartgetID is empty, cannot create DM channel")
	}

	channel, err := s.session.UserChannelCreate(targetID)
	if err != nil {
		return fmt.Errorf("SendNotification | error creating channel for %s: %w", targetID, err)
	}

	s.msgInvite(targetID, set, channel)

	return nil
}

func (s *DiscordSender) GetPlatformMessenagerName() string {
	return "discord"
}

func (s *DiscordSender) FindContactOfParticipant(ctx context.Context, p sender.Participant) (sender.Participant, error) {
	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	request := entity.ParticipantGetRequest{
		GamerTag:           p.GameNickname,
		MessenagerPlatform: s.GetPlatformMessenagerName(),
	}

	response, err := s.participantUC.GetParticipant(request)
	if err == nil {
		return sender.Participant{
			MessenagerID:    response.MessengerPlatformId,
			MessenagerLogin: p.MessenagerLogin,
			MessenagerName:  s.GetPlatformMessenagerName(),
			GameNickname:    p.GameNickname,
			GameID:          p.GameID,
			Locales:         []string{response.Locale},
		}, nil
	}
	log.Printf("db | participant %s not found, searching in Discord...", p.GameNickname)

	cleanNickname := s.cleanDiscordLogin(p.MessenagerLogin)
	var messengerID string
	currentLocale := "default"

	if p.MessenagerLogin == "" || p.MessenagerLogin == "N/D" {
		if s.debugMode {
			log.Printf("findContact | %s has no login, using debug mock", p.GameNickname)
			messengerID = "000000000000000000"
			cleanNickname = "N/D"
			currentLocale = s.cfg.rolesIdList.Ru
		} else {
			return sender.Participant{}, fmt.Errorf("findContact | member %s not founded in guild (server)\n", cleanNickname)
		}
	} else {
		members, err := s.session.GuildMembersSearch(s.cfg.guildID, cleanNickname, 1)
		if err != nil || len(members) != 1 {
			if s.debugMode {
				messengerID = "000000000000000000"
				currentLocale = s.cfg.rolesIdList.Ru
			} else {
				return sender.Participant{}, fmt.Errorf("findContact | member %s not founded in guild (server)\n", cleanNickname)
			}
		} else {
			targetMember := members[0]
			messengerID = targetMember.User.ID

			for _, roleId := range targetMember.Roles {
				// TODO: Change reconize locale in future (More languages)
				if roleId == s.cfg.rolesIdList.Ru {
					currentLocale = roleId
				}
			}
		}
	}

	addRequest := entity.ParticipantAddRequest{
		GamerTag:               p.GameNickname,
		MessengerPlatform:      s.GetPlatformMessenagerName(),
		MessengerPlatformId:    messengerID,
		MessengerPlatformLogin: cleanNickname,
		IsFound:                true,
		UpdatedAt:              time.Now(),
		Locale:                 currentLocale,
	}

	_, err = s.participantUC.AddParticipant(addRequest)
	if err != nil {
		log.Printf("db | failed to save participant %v: %v", cleanNickname, err)
	} else {
		log.Printf("db | successfully saved participant %v", cleanNickname)
	}

	return sender.Participant{
		MessenagerID:    addRequest.MessengerPlatformId,
		MessenagerLogin: addRequest.MessengerPlatformLogin,
		MessenagerName:  s.GetPlatformMessenagerName(),
		GameNickname:    p.GameNickname,
		GameID:          p.GameID,
		Locales:         []string{addRequest.Locale},
	}, nil
}

func (s *DiscordSender) cleanDiscordLogin(login string) string {
	res := strings.ReplaceAll(login, "@", "")
	if strings.Contains(res, "#") {
		return strings.Split(res, "#")[0]
	}
	return res
}

func (dh *discordHandler) Process(s *discordgo.Session) {
	dh.mutex.Lock()

	if dh.cancel != nil {
		dh.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	dh.cancel = cancel
	dh.mutex.Unlock()

	defer func() {
		cancel()
		dh.mutex.Lock()
		dh.cancel = nil
		dh.mutex.Unlock()
	}()

	if err := dh.SendingMessages(ctx); err != nil {
		log.Printf("SendingMessages stopped or failed: %v", err)
	}
}

func (dh *discordHandler) getAdapter() (sender.NotificationData, error) {
	switch dh.tournamentPlatform {
	case "startgg":
		client, err := auth.GetClientStartgg()
		if err != nil {
			return nil, err
		}

		contacts, err := senderUC.LoadCSV("contacts.json")

		return senderUC.StartggSetAdapter{
			FullSlug:  dh.slug,
			Client:    client,
			DebugMode: dh.debugMode,
			Contacts:  contacts,
		}, nil
	case "challonge":
		client, err := auth.GetClientChallonge()
		if err != nil {
			return nil, err
		}

		return senderUC.ChallongeMatchAdapter{
			TournamentSlug: dh.slug,
			Client:         client,
			DebugMode:      dh.debugMode,
		}, nil
	default:
		return nil, fmt.Errorf("getAdapter | Can't get adapter for platform called: %s", dh.tournamentPlatform)
	}
}

func (dh *discordHandler) SendingMessages(ctx context.Context) error {
	adapter, err := dh.getAdapter()
	if err != nil {
		return err
	}
	dh.ns.Data = adapter
	if err := dh.ns.Process(ctx); err != nil {
		return fmt.Errorf("sendingMessages | process failed: %w", err)
	}

	return nil
}
