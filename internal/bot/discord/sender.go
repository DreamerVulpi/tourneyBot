package discord

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"

	"time"

	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
	"github.com/dreamervulpi/tourneyBot/internal/db/usecase"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	"github.com/dreamervulpi/tourneyBot/locale"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type DiscordSender struct {
	session       *discordgo.Session
	cfg           params
	participantUC usecase.Participant
	// TODO: Change startgg on universal
	startgg   strtgg
	debugMode bool
}

func (s *DiscordSender) prepareSetData(set sender.SetData, local locale.Lang, link string) (*discordgo.MessageEmbed, error) {
	format := s.cfg.rulesMatches.StandardFormat
	if s.startgg.finalBracketId == set.PhaseGroupId {
		if s.startgg.minRoundNumA <= set.RoundNum && set.RoundNum <= s.startgg.minRoundNumB || s.startgg.maxRoundNumA <= set.RoundNum && set.RoundNum <= s.startgg.maxRoundNumB {
			format = s.cfg.rulesMatches.FinalsFormat
		}
	}

	embedColor := ColorDefault
	if len(set.StreamSourse) > 0 {
		embedColor = ColorStream
	}

	var message *discordgo.MessageEmbed
	gameID := set.Opponent.GameID
	if len(gameID) == 0 {
		gameID = local.ErrorMessage.NoData
	}
	rawID := set.Opponent.MessenagerID
	login := set.Opponent.MessenagerLogin
	var discordDisplay string
	if len(rawID) == 0 && len(login) == 0 {
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
	log.Printf("DiscordDisplay: %v | Set: %v", discordDisplay, link)

	if len(set.StreamSourse) == 0 {
		fields := []*discordgo.MessageEmbedField{
			{Name: local.InviteMessage.MessageHeader},
			{Name: local.InviteMessage.Nickname, Value: fmt.Sprintf("```%v```", set.Opponent.GameNickname), Inline: true},
			{Name: local.InviteMessage.GameID, Value: fmt.Sprintf("```%v```", gameID), Inline: true},
			{Name: local.InviteMessage.Discord, Value: discordDisplay, Inline: true},

			{Name: local.InviteMessage.CheckIn, Value: link},
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
			{Name: local.StreamLobbyMessage.MessageHeader, Value: link},
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

func (s *DiscordSender) msgInvite(set sender.SetData, channel *discordgo.Channel, link string, roleId string) {
	var local locale.Lang
	if roleId == s.cfg.rolesIdList.Ru {
		local = locale.Ru
	} else {
		local = locale.En
	}

	message, err := s.prepareSetData(set, local, link)
	if err != nil {
		log.Printf("msgInvite | error sended DM: %v\n", err.Error())
		failMsg := []*discordgo.MessageEmbedField{
			{Name: local.LogMessage.FailMsgHeader},
			{Name: fmt.Sprintf(local.LogMessage.FailedSentMsg, set.Recipient.MessenagerLogin), Value: "\u200B"},
		}
		failedLog := s.templateEmbedMsg(fmt.Sprintf(local.LogMessage.Title, set.TournamentName), failMsg, ColorError)
		if _, err := s.session.ChannelMessageSendEmbed(s.cfg.debugChannelID, failedLog); err != nil {
			log.Printf("msgInvite | error sending to debugChannel: %v\n", err.Error())
		}
		return
	}

	_, err = s.session.ChannelMessageSendEmbed(channel.ID, message)
	if err != nil {
		log.Printf("msgInvite | error sended DM: %v\n", err.Error())
		failMsg := []*discordgo.MessageEmbedField{
			{Name: local.LogMessage.FailMsgHeader},
			{Name: fmt.Sprintf(local.LogMessage.FailedSentMsg, set.Recipient.MessenagerLogin), Value: "\u200B"},
		}
		failedLog := s.templateEmbedMsg(fmt.Sprintf(local.LogMessage.Title, set.TournamentName), failMsg, ColorError)
		if _, err := s.session.ChannelMessageSendEmbed(s.cfg.debugChannelID, failedLog); err != nil {
			log.Printf("msgInvite | error sending to debugChannel: %v\n", err.Error())
		}
	}

	successMsg := []*discordgo.MessageEmbedField{
		{Name: local.LogMessage.SuccessfulMsgHeader},
		{Name: fmt.Sprintf(local.LogMessage.SuccesfullSendedMsg, set.Recipient.MessenagerLogin), Value: "\u200B"},
	}
	successLog := s.templateEmbedMsg(fmt.Sprintf(local.LogMessage.Title, set.TournamentName), successMsg, ColorSuccess)
	if _, err := s.session.ChannelMessageSendEmbed(s.cfg.debugChannelID, successLog); err != nil {
		log.Printf("msgInvite | error sending to debugChannel: %v\n", err.Error())
		return
	}
	log.Printf("msgInvite | success sended to DM")
}

func (s *DiscordSender) SendNotification(ctx context.Context, targetID string, data sender.SetData) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	channel, err := s.session.UserChannelCreate(data.Recipient.MessenagerID)
	if err != nil {
		return fmt.Errorf("discord: error creating channel for %s: %w", data.Recipient.MessenagerID, err)
	}

	if len(data.Recipient.Locales) != 0 {
		s.msgInvite(data, channel, data.FullInviteLink, data.Recipient.Locales[0])
	} else {
		s.msgInvite(data, channel, data.FullInviteLink, "default")
	}

	return nil
}

func (s *DiscordSender) GetPlatformMessenagerName() string {
	return "discord"
}

func (s *DiscordSender) FindContactOfParticipant(ctx context.Context, p sender.Participant) (sender.Participant, error) {
	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	if p.MessenagerLogin == "" || p.MessenagerLogin == "N/D" {
		return sender.Participant{}, fmt.Errorf("findContact | empty platform login for %v\n", p.GameNickname)
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
	locale := "default"

	members, err := s.session.GuildMembersSearch(s.cfg.guildID, cleanNickname, 1)
	if err != nil || len(members) != 1 {
		if s.debugMode {
			log.Printf("search: %s not found, using mock data (debug)\n", cleanNickname)
			messengerID = "000000000000000000"
			locale = s.cfg.rolesIdList.Ru
		} else {
			return sender.Participant{}, fmt.Errorf("findContact | member %s not founded in guild (server)\n", cleanNickname)
		}
	} else {
		targetMember := members[0]
		messengerID = targetMember.User.ID

		for _, roleId := range targetMember.Roles {
			// TODO: Нужно поменять логику локали. Если их несколько??
			if roleId == s.cfg.rolesIdList.Ru {
				locale = roleId
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
		Locale:                 locale,
	}

	_, err = s.participantUC.AddParticipant(addRequest)
	if err != nil {
		log.Printf("db | failed to save participant %v: %v", cleanNickname, err)
		return sender.Participant{}, err
	}

	log.Printf("db | successfully saved participant %v", cleanNickname)

	return sender.Participant{
		MessenagerID:    addRequest.MessengerPlatformId,
		MessenagerLogin: addRequest.GamerTag,
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

// REFACTOR:
func (s *discordHandler) checkContact(participant startgg.Participant) sender.Participant {
	p := sender.Participant{
		MessenagerLogin: "N/D",
		GameID:          "N/D",
		GameNickname:    "N/D",
	}

	// first participant from team (solo)
	src := participant
	p.GameNickname = src.GamerTag

	// search discord login in profile startgg
	if len(src.User.Authorizations) > 0 {
		p.MessenagerLogin = src.User.Authorizations[0].Discord
	}

	// if empty then check local file json
	if p.MessenagerLogin == "N/D" || p.MessenagerLogin == "" {
		if val, ok := s.contacts.contacts[src.GamerTag]; ok {
			p.MessenagerLogin = val.MessenagerLogin
		}
	}

	// get game ID from startgg
	apiID := ""
	switch s.cfg.tournament.Game.Name {
	case "tekken":
		apiID = src.ConnectedAccounts.Tekken.TekkenID
	case "sf6":
		apiID = src.ConnectedAccounts.SF6.GameID
	}

	// if empty then check local file json
	if apiID != "" {
		p.GameID = apiID
	} else {
		if val, ok := s.contacts.contacts[strings.ToLower(src.GamerTag)]; ok {
			p.GameID = val.GameID
		}
	}

	return p
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

	// TODO: Must be stoppable
	if err := dh.SendingMessages(ctx); err != nil {
		log.Printf("SendingMessages stopped or failed: %v", err)
	}
}

// TODO: On future
// func (dh *commandHandler) checkPhaseGroup(phaseGroupId int64, sets []startgg.Nodes) error {
// 	var max, minIndex, maxIndex int
// 	min := sets[0].Round
// 	for index, set := range sets {
// 		if min > set.Round {
// 			min = set.Round
// 			minIndex = index
// 		}
// 		if max < set.Round {
// 			max = set.Round
// 			maxIndex = index
// 		}
// 	}
// 	if sets[maxIndex].FullRoundText == "Grand Final" && sets[minIndex].FullRoundText == "Losers Final" {
// 		log.Printf("Finded final bracket! -> %v , %v\n", min, max)
// 		dh.startgg.minRoundNumA = min
// 		dh.startgg.maxRoundNumB = max
// 		dh.startgg.minRoundNumB = dh.startgg.minRoundNumA + 2
// 		dh.startgg.maxRoundNumA = dh.startgg.maxRoundNumB - 3
// 		dh.startgg.finalBracketId = phaseGroupId
// 		return nil
// 	} else {
// 		return errors.New("not final bracket")
// 	}
// }

func (dh *discordHandler) SendingMessages(ctx context.Context) error {
	if dh.auth == nil {
		return errors.New("sendingMessages: auth client is not initialized - check bot.Start parameters")
	}

	var testUser sender.Participant
	if dh.debugMode {
		var err error
		me, err := dh.auth.GetDiscordMe(ctx)
		if err != nil {
			return fmt.Errorf("debug setup failed: %w", err)
		}

		testUser = sender.Participant{
			MessenagerID:    me.ID,
			MessenagerLogin: me.Username,
			Locales:         []string{"ru"},
		}
		log.Printf("My PlatformID: %v\n", testUser.MessenagerID)
	}

	// REFACTOR: Change to interface NotificationData
	// Get pages with state: Not started
	states := []int{1}
	if dh.debugMode {
		states = []int{1, 2, 3}
	}

	tournament, err := dh.startgg.client.GetTournament(strings.Replace(strings.SplitAfter(dh.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	phaseGroups, err := dh.startgg.client.GetListGroups(dh.slug)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, phaseGroupId := range phaseGroups {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		total, err := dh.startgg.client.GetPagesCount(phaseGroupId.Id, states)
		if err != nil || total == 0 {
			continue
		}

		var pages int
		if total <= 60 {
			pages = 1
		} else {
			pages = int(math.Ceil(float64(total) / 60.0))
		}

		log.Printf("%v | %v | Pages: %v\n", phaseGroupId, total, pages)

		for i := 0; i < pages; i++ {
			log.Printf("%v | Page #%v\n", phaseGroupId, i)
			sets, err := dh.startgg.client.GetSets(phaseGroupId.Id, i+1, 60, states)
			if err != nil {
				log.Printf("Error getting sets: %v", err)
				continue
			}

			for _, set := range sets {
				if ctx.Err() != nil {
					break
				}

				request := entity.SentSetCheckRequest{SetId: set.Id}
				alreadySent, err := dh.sentSetUC.IsExists(request)
				if err != nil {
					log.Printf("DB error checking set %v: %v", set.Id, err)
				}

				if alreadySent.State && !dh.debugMode {
					continue
				}

				wg.Add(1)
				go func(ctx context.Context, set startgg.Nodes, testUser sender.Participant) {
					defer wg.Done()

					if ctx.Err() != nil {
						return
					}

					// TODO: ADD CHECK 0 COUNT
					if len(set.Slots[0].Entrant.Participants) == 0 {
						log.Printf("No contact data from startgg")
					}
					if len(set.Slots[1].Entrant.Participants) == 0 {
						log.Printf("No contact data from startgg")
					}

					// discord contact check
					p1 := dh.checkContact(set.Slots[0].Entrant.Participants[0])
					p2 := dh.checkContact(set.Slots[1].Entrant.Participants[0])

					if ctx.Err() != nil {
						return
					}

					contactP1, err := dh.msgSender.FindContactOfParticipant(ctx, p1)
					if err != nil {
						log.Printf("FindContact P1 Error (%s): %v\n", p1.MessenagerLogin, err)
					}
					//  else {
					// 	p1.MessenagerID = contactP1.MessenagerID
					// 	p1.Locales = contactP1.Locales
					// }

					if ctx.Err() != nil {
						return
					}

					contactP2, err := dh.msgSender.FindContactOfParticipant(ctx, p2)
					if err != nil {
						log.Printf("FindContact P2 Error (%s): %v\n", p2.MessenagerLogin, err)
					}
					//  else {
					// 	p2.MessenagerID = contactP2.MessenagerID
					// 	p2.Locales = contactP2.Locales
					// }

					if ctx.Err() != nil {
						return
					}

					toPlayer1 := sender.SetData{
						TournamentName: tournament.Name,
						SetID:          set.Id,
						StreamName:     set.Stream.StreamName,
						StreamSourse:   set.Stream.StreamSource,
						RoundNum:       set.Round,
						PhaseGroupId:   phaseGroupId.Id,
						Recipient:      contactP1, // For player1
						Opponent:       contactP2,
						// sender.Participant{
						// 	MessenagerID:    contactP2.MessenagerID, // To player2
						// 	MessenagerLogin: p2.MessenagerLogin,
						// 	GameNickname:    set.Slots[1].Entrant.Participants[0].GamerTag,
						// 	GameID:          p2.GameID,
						// },
						FullInviteLink: fmt.Sprint("https://www.start.gg/", dh.slug, "/set/", set.Id),
					}
					toPlayer2 := sender.SetData{
						TournamentName: tournament.Name,
						SetID:          set.Id,
						StreamName:     set.Stream.StreamName,
						StreamSourse:   set.Stream.StreamSource,
						RoundNum:       set.Round,
						PhaseGroupId:   phaseGroupId.Id,
						Recipient:      contactP2, // For player2
						Opponent:       contactP2,
						// sender.Participant{
						// 	MessenagerID:    contactP1.MessenagerID, // To player1
						// 	MessenagerLogin: p1.MessenagerLogin,
						// 	GameNickname:    set.Slots[0].Entrant.Participants[0].GamerTag,
						// 	GameID:          p1.GameID,
						// },
						FullInviteLink: fmt.Sprint("https://www.start.gg/", dh.slug, "/set/", set.Id),
					}

					if dh.debugMode {
						toPlayer1.Recipient = testUser
						toPlayer2.Recipient = testUser
					}

					if dh.debugMode || p1.MessenagerLogin != "N/D" {
						if err := dh.msgSender.SendNotification(ctx, toPlayer1.Recipient.MessenagerID, toPlayer1); err != nil {
							// TODO: Sent error to debugChannel and solution
							log.Printf("Error sending to p1 (%s): %v", toPlayer1.Recipient.MessenagerID, err)
						}

						select {
						case <-ctx.Done():
							return
						case <-time.After(1 * time.Second):
						}
					}

					if dh.debugMode || p2.MessenagerLogin != "N/D" {
						if err := dh.msgSender.SendNotification(ctx, toPlayer2.Recipient.MessenagerID, toPlayer2); err != nil {
							// TODO: Sent error to debugChannel and solution
							log.Printf("Error sending to p2 (%s): %v", toPlayer2.Recipient.MessenagerID, err)
						}
						select {
						case <-ctx.Done():
							return
						case <-time.After(1 * time.Second):
						}
					}
					request := entity.SentSetAddRequest{
						SetId: set.Id,
						// TODO: Изменить на универсальную платформу турниров
						TournamentPlatform: "startgg",
						MessengerPlatform:  dh.msgSender.GetPlatformMessenagerName(),
						TournamentSlug:     dh.slug,
						SentAt:             time.Now(),
					}
					_, err = dh.sentSetUC.AddSentSet(request)
					if err != nil {
						log.Printf("DB error saving sent set %v: %v", set.Id, err)
					}

				}(ctx, set, testUser)
			}
			log.Printf("Checked phaseGroup(%v)", phaseGroupId)
		}
	}

	wg.Wait()
	return nil
}
