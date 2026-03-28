package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	entityDB "github.com/dreamervulpi/tourneyBot/internal/entity/db"
	entitySender "github.com/dreamervulpi/tourneyBot/internal/entity/sender"
	usecaseDB "github.com/dreamervulpi/tourneyBot/internal/usecase/db"
	senderUC "github.com/dreamervulpi/tourneyBot/internal/usecase/sender"
)

type DiscordSender struct {
	session                *discordgo.Session
	cfg                    params
	participantUC          usecaseDB.Participant
	adminIDTourneyPlatform string
	debugMode              bool
}

func (dh *DiscordHandler) Process(s *discordgo.Session) {
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

func (dh *DiscordHandler) SendingMessages(ctx context.Context) error {
	adapter, err := dh.GetAdapter()
	if err != nil {
		return err
	}
	dh.ns.Data = adapter
	if err := dh.ns.Process(ctx); err != nil {
		return fmt.Errorf("sendingMessages | process failed: %w", err)
	}

	return nil
}

func (s *DiscordSender) SendNotification(ctx context.Context, targetID string, set entitySender.SetData) error {
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

	message, local, recipient := s.msgInvite(targetID, set, channel)

	_, err = s.session.ChannelMessageSendEmbed(channel.ID, message)
	if err != nil {
		log.Printf("SendNotification | error sended DM: %v\n", err.Error())
		s.logMsgToDiscord(false, err.Error(), set, local, recipient.GameNickname)
		return err
	}

	s.logMsgToDiscord(true, "", set, local, recipient.GameNickname)
	return nil
}

func (s *DiscordSender) FindContactOfParticipant(ctx context.Context, p entitySender.Participant) (entitySender.Participant, error) {
	if err := ctx.Err(); err != nil {
		return entitySender.Participant{}, err
	}

	request := entityDB.ParticipantGetRequest{
		GamerTag:           p.GameNickname,
		MessenagerPlatform: s.GetPlatformMessenagerName(),
	}

	response, err := s.participantUC.GetParticipant(request)
	if err == nil {
		return entitySender.Participant{
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
			return entitySender.Participant{}, fmt.Errorf("findContact | member %s not founded in guild (server)\n", cleanNickname)
		}
	} else {
		members, err := s.session.GuildMembersSearch(s.cfg.guildID, cleanNickname, 1)
		if err != nil || len(members) != 1 {
			if s.debugMode {
				messengerID = "000000000000000000"
				currentLocale = s.cfg.rolesIdList.Ru
			} else {
				return entitySender.Participant{}, fmt.Errorf("findContact | member %s not founded in guild (server)\n", cleanNickname)
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

	addRequest := entityDB.ParticipantAddRequest{
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

	return entitySender.Participant{
		MessenagerID:    addRequest.MessengerPlatformId,
		MessenagerLogin: addRequest.MessengerPlatformLogin,
		MessenagerName:  s.GetPlatformMessenagerName(),
		GameNickname:    p.GameNickname,
		GameID:          p.GameID,
		Locales:         []string{addRequest.Locale},
	}, nil
}

func (s *DiscordSender) GetPlatformMessenagerName() string {
	return "discord"
}

func (s *DiscordSender) cleanDiscordLogin(login string) string {
	res := strings.ReplaceAll(login, "@", "")
	if strings.Contains(res, "#") {
		return strings.Split(res, "#")[0]
	}
	return res
}

func (dh *DiscordHandler) GetAdapter() (entitySender.NotificationData, error) {
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

		// TODO: Load contacts from file for challonge

		return senderUC.ChallongeMatchAdapter{
			TournamentSlug: dh.slug,
			Client:         client,
			DebugMode:      dh.debugMode,
		}, nil
	default:
		return nil, fmt.Errorf("getAdapter | Can't get adapter for platform called: %s", dh.tournamentPlatform)
	}
}
