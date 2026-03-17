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
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	"github.com/dreamervulpi/tourneyBot/locale"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type DiscordSender struct {
	session *discordgo.Session
	config  params
	// TODO: Change startgg on universal
	startgg        strtgg
	debugChannelID string
}

func (s *DiscordSender) prepareSetData(set sender.SetData, local locale.Lang, link string) (*discordgo.MessageEmbed, error) {
	format := s.config.rulesMatches.StandardFormat
	if s.startgg.finalBracketId == set.PhaseGroupId {
		if s.startgg.minRoundNumA <= set.RoundNum && set.RoundNum <= s.startgg.minRoundNumB || s.startgg.maxRoundNumA <= set.RoundNum && set.RoundNum <= s.startgg.maxRoundNumB {
			format = s.config.rulesMatches.FinalsFormat
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
			{Name: fmt.Sprintf(local.InviteMessage.Warning, s.config.rulesMatches.Waiting), Value: "\u200B"},

			{Name: local.InviteMessage.SettingsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.InviteMessage.Stage, Value: s.fieldStage(local), Inline: true},
			{Name: local.InviteMessage.Rounds, Value: fmt.Sprintf("%v", s.config.rulesMatches.Rounds), Inline: true},
			{Name: local.InviteMessage.Duration, Value: fmt.Sprintf(local.InviteMessage.DurationCount, s.config.rulesMatches.Duration), Inline: true},
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
			{Name: fmt.Sprintf(local.StreamLobbyMessage.Warning, s.config.rulesMatches.Waiting), Value: "\u200B"},

			{Name: local.StreamLobbyMessage.ParamsHeader},
			{Name: local.InviteMessage.StandardFormat, Value: fmt.Sprintf(local.InviteMessage.FT, format) + fmt.Sprintf(local.InviteMessage.FormatDescription, format), Inline: true},
			{Name: local.StreamLobbyMessage.Area, Value: s.fieldArea(local), Inline: true},
			{Name: local.StreamLobbyMessage.Language, Value: s.fieldLanguage(local), Inline: true},
			{Name: local.StreamLobbyMessage.TypeConnection, Value: s.fieldConnection(local), Inline: true},
			{Name: local.StreamLobbyMessage.Crossplatform, Value: s.fieldCrossplay(local), Inline: true},
			{Name: local.StreamLobbyMessage.Passcode, Value: fmt.Sprintf(local.StreamLobbyMessage.PasscodeTemplate, s.config.streamLobby.Passcode), Inline: true},
		}
		message = s.templateEmbedMsg(fmt.Sprintf(local.StreamLobbyMessage.Title, set.TournamentName), fields, embedColor)
		message.Description = local.StreamLobbyMessage.Description
	}
	return message, nil
}

func (s *DiscordSender) msgInvite(set sender.SetData, channel *discordgo.Channel, link string, roleId string) {
	var local locale.Lang
	if roleId == s.config.rolesIdList.Ru {
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
		if _, err := s.session.ChannelMessageSendEmbed(s.config.debugChannelID, failedLog); err != nil {
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
		if _, err := s.session.ChannelMessageSendEmbed(s.config.debugChannelID, failedLog); err != nil {
			log.Printf("msgInvite | error sending to debugChannel: %v\n", err.Error())
		}
	}

	successMsg := []*discordgo.MessageEmbedField{
		{Name: local.LogMessage.SuccessfulMsgHeader},
		{Name: fmt.Sprintf(local.LogMessage.SuccesfullSendedMsg, set.Recipient.MessenagerLogin), Value: "\u200B"},
	}
	successLog := s.templateEmbedMsg(fmt.Sprintf(local.LogMessage.Title, set.TournamentName), successMsg, ColorSuccess)
	if _, err := s.session.ChannelMessageSendEmbed(s.config.debugChannelID, successLog); err != nil {
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

func (dh *discordHandler) searchContactDiscord(ctx context.Context, s *discordgo.Session, platformNickname string, gameNickname string) (sender.Participant, error) {
	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	if platformNickname == "" || platformNickname == "N/D" {
		return sender.Participant{}, fmt.Errorf("searchContactDiscord: empty platformNickname %v", platformNickname)
	}

	cleanNickname := strings.Split(platformNickname, "#")[0]

	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	request := entity.ParticipantGetRequest{
		GamerTag:           gameNickname,
		MessenagerPlatform: dh.msgSender.GetPlatformMessenagerName(),
	}

	response, err := dh.participantUC.GetParticipant(request)
	if err != nil {
		log.Printf("db | not finded player %v from %v in database", request.GamerTag, request.MessenagerPlatform)

		if err := ctx.Err(); err != nil {
			return sender.Participant{}, err
		}

		var mpi string
		locale := "N/D"

		members, err := s.GuildMembersSearch(dh.cfg.guildID, cleanNickname, 1)
		if (err != nil || len(members) != 1) && !dh.debugMode {
			return sender.Participant{}, fmt.Errorf("searchContactDiscord: player not finded %v", cleanNickname)
		}

		if (err != nil || len(members) != 1) && dh.debugMode {
			log.Printf("searchContactDiscord: member %v not found on DiscordServer, using mock data", cleanNickname)
			mpi = "000000000000000000"
			locale = dh.cfg.rolesIdList.Ru
		} else {
			targetMember := members[0]
			// Get list rolesId including in locale (en is default)
			mpi = targetMember.User.ID
			for _, roleId := range targetMember.Roles {
				// TODO: Нужно поменять логику локали. Если их несколько??
				if roleId == dh.cfg.rolesIdList.Ru {
					locale = roleId
				}
			}
		}

		if len(mpi) == 0 {
			mpi = "N/D"
		}

		addRequest := entity.ParticipantAddRequest{
			GamerTag:               gameNickname,
			MessengerPlatform:      dh.msgSender.GetPlatformMessenagerName(),
			MessengerPlatformId:    mpi,
			MessengerPlatformLogin: cleanNickname,
			IsFound:                true,
			UpdatedAt:              time.Now(),
			Locale:                 locale,
		}

		if err := ctx.Err(); err != nil {
			return sender.Participant{}, err
		}

		_, err = dh.participantUC.AddParticipant(addRequest)
		if err != nil {
			log.Printf("db | failed to save participant %v: %v", cleanNickname, err)
			return sender.Participant{}, err
		}

		log.Printf("db | successfully saved participant %v", cleanNickname)
		return sender.Participant{
			MessenagerID:    addRequest.MessengerPlatformId,
			MessenagerLogin: addRequest.GamerTag,
			Locales:         []string{addRequest.Locale},
		}, nil

	}

	return sender.Participant{
		MessenagerID:    response.MessengerPlatformId,
		MessenagerLogin: response.GamerTag,
		Locales:         []string{response.Locale},
	}, nil
}

func (s *discordHandler) checkContact(participants []startgg.Participants) sender.Participant {
	p := sender.Participant{
		MessenagerLogin: "N/D",
		GameID:          "N/D",
		GameNickname:    "N/D",
	}

	if len(participants) == 0 {
		return sender.Participant{}
	}

	// first participant from team (solo)
	src := participants[0]
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

	if err := dh.SendingMessages(ctx, s); err != nil {
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

// REFACTOR: Вынести в отдельно
func (dh *discordHandler) SendingMessages(ctx context.Context, s *discordgo.Session) error {
	if dh.auth == nil {
		return errors.New("sendingMessages: auth client is not initialized - check bot.Start parameters")
	}

	tournament, err := dh.startgg.client.GetTournament(strings.Replace(strings.SplitAfter(dh.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	phaseGroups, err := dh.startgg.client.GetListGroups(dh.slug)
	if err != nil {
		return err
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

	// Get pages with state: Not started
	states := []int{1}
	if dh.debugMode {
		states = []int{1, 2, 3}
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

					// discord contact check
					p1 := dh.checkContact(set.Slots[0].Entrant.Participants)
					p2 := dh.checkContact(set.Slots[1].Entrant.Participants)

					if ctx.Err() != nil {
						return
					}

					contactP1, err := dh.searchContactDiscord(ctx, s, p1.MessenagerLogin, p1.GameNickname)
					if err != nil {
						log.Printf("sending message | Not finded member in discord (%v)\n", p1.MessenagerLogin)
					} else {
						p1.MessenagerID = contactP1.MessenagerID
						p1.Locales = contactP1.Locales
					}

					if ctx.Err() != nil {
						return
					}

					contactP2, err := dh.searchContactDiscord(ctx, s, p2.MessenagerLogin, p2.GameNickname)
					if err != nil {
						log.Printf("sending message | Not finded member in discord (%v)\n", p2.MessenagerLogin)
					} else {
						p2.MessenagerID = contactP2.MessenagerID
						p2.Locales = contactP2.Locales
					}

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
						Opponent: sender.Participant{
							MessenagerID:    contactP2.MessenagerID, // To player2
							MessenagerLogin: p2.MessenagerLogin,
							GameNickname:    set.Slots[1].Entrant.Participants[0].GamerTag,
							GameID:          p2.GameID,
						},
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
						Opponent: sender.Participant{
							MessenagerID:    contactP1.MessenagerID, // To player1
							MessenagerLogin: p1.MessenagerLogin,
							GameNickname:    set.Slots[0].Entrant.Participants[0].GamerTag,
							GameID:          p1.GameID,
						},
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
