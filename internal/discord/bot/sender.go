package bot

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
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type PlayerData struct {
	tournament   string
	setID        int64
	streamName   string
	streamSourse string
	roundNum     int
	phaseGroupId int64
	recipient    sender.Participant
	opponent     sender.Participant
	IsTest       bool
}

func (ch *commandHandler) searchContactDiscord(ctx context.Context, s *discordgo.Session, nickname string) (sender.Participant, error) {
	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	if nickname == "" || nickname == "N/D" {
		return sender.Participant{}, fmt.Errorf("searchContactDiscord: empty nickname %v", nickname)
	}

	cleanName := strings.Split(nickname, "#")[0]

	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	members, err := s.GuildMembersSearch(ch.cfg.guildID, cleanName, 1)
	if err != nil {
		return sender.Participant{}, err
	}

	if err := ctx.Err(); err != nil {
		return sender.Participant{}, err
	}

	if len(members) != 1 {
		return sender.Participant{}, fmt.Errorf("searchContactDiscord: player not finded %v", cleanName)
	}

	targetMember := members[0]
	// Get list rolesId including in locale (en is default)
	roles := []string{}
	for _, roleId := range targetMember.Roles {
		if roleId == ch.cfg.rolesIdList.Ru {
			roles = append(roles, roleId)
		}
	}

	return sender.Participant{
		DiscordID:    strings.Split(targetMember.User.ID, "#")[0],
		DiscordLogin: nickname,
		Locales:      roles,
	}, nil
}

func (ch *commandHandler) checkContact(participants []startgg.Participants) sender.Participant {
	p := sender.Participant{
		DiscordLogin: "N/D",
		GameID:       "N/D",
		GameNickname: "N/D",
	}

	if len(participants) == 0 {
		return sender.Participant{}
	}

	// first participant from team (solo)
	src := participants[0]
	p.GameNickname = src.GamerTag

	// search discord login in profile startgg
	if len(src.User.Authorizations) > 0 {
		p.DiscordLogin = src.User.Authorizations[0].Discord
	}

	// if empty then check local file json
	if p.DiscordLogin == "N/D" || p.DiscordLogin == "" {
		if val, ok := ch.discord.contacts[src.GamerTag]; ok {
			p.DiscordLogin = val.DiscordLogin
		}
	}

	// get game ID from startgg
	apiID := ""
	switch ch.cfg.tournament.Game.Name {
	case "tekken":
		apiID = src.ConnectedAccounts.Tekken.TekkenID
	case "sf6":
		apiID = src.ConnectedAccounts.SF6.GameID
	}

	// if empty then check local file json
	if apiID != "" {
		p.GameID = apiID
	} else {
		if val, ok := ch.discord.contacts[strings.ToLower(src.GamerTag)]; ok {
			p.GameID = val.GameID
		}
	}

	return p
}

func (ch *commandHandler) sendMsg(ctx context.Context, s *discordgo.Session, player PlayerData) {
	if err := ctx.Err(); err != nil {
		return
	}

	if player.recipient.DiscordID == "" {
		log.Printf("Skip sending: empty Discord ID for player %v", player.opponent.GameNickname)
		return
	}

	channel, err := s.UserChannelCreate(player.recipient.DiscordID)
	if err != nil {
		log.Println("error creating channel:", err)
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	link := fmt.Sprint("https://www.start.gg/", ch.slug, "/set/", player.setID)

	if len(player.recipient.Locales) != 0 {
		ch.msgInvite(s, player, channel, link, player.recipient.Locales[0])
	} else {
		ch.msgInvite(s, player, channel, link, "default")
	}
}

func (ch *commandHandler) Process(s *discordgo.Session) {
	ch.mu.Lock()

	if ch.cancelFunc != nil {
		ch.cancelFunc()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch.cancelFunc = cancel
	ch.mu.Unlock()

	defer func() {
		cancel()
		ch.mu.Lock()
		ch.cancelFunc = nil
		ch.mu.Unlock()
	}()

	if err := ch.SendingMessages(ctx, s); err != nil {
		log.Printf("SendingMessages stopped or failed: %v", err)
	}
}

// func (ch *commandHandler) checkPhaseGroup(phaseGroupId int64, sets []startgg.Nodes) error {
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
// 		ch.startgg.minRoundNumA = min
// 		ch.startgg.maxRoundNumB = max
// 		ch.startgg.minRoundNumB = ch.startgg.minRoundNumA + 2
// 		ch.startgg.maxRoundNumA = ch.startgg.maxRoundNumB - 3
// 		ch.startgg.finalBracketId = phaseGroupId
// 		return nil
// 	} else {
// 		return errors.New("not final bracket")
// 	}
// }

func (ch *commandHandler) SendingMessages(ctx context.Context, s *discordgo.Session) error {
	if ch.auth == nil {
		return errors.New("sendingMessages: auth client is not initialized - check bot.Start parameters")
	}
	var wg sync.WaitGroup

	tournament, err := ch.startgg.client.GetTournament(strings.Replace(strings.SplitAfter(ch.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	phaseGroups, err := ch.startgg.client.GetListGroups(ch.slug)
	if err != nil {
		return err
	}

	var testUser sender.Participant
	if ch.debugMode {
		var err error
		me, err := ch.auth.GetDiscordMe(ctx)
		if err != nil {
			return fmt.Errorf("debug setup failed: %w", err)
		}

		testUser = sender.Participant{
			DiscordID:    me.ID,
			DiscordLogin: me.Username,
			Locales:      []string{"ru"},
		}
		log.Printf("My DiscordID: %v\n", testUser.DiscordID)
	}

	// Get pages with state: Not started
	states := []int{1}
	if ch.debugMode {
		states = []int{1, 2, 3}
	}

	for _, phaseGroupId := range phaseGroups {
		// Проверка контекста перед началом обработки новой группы
		if ctx.Err() != nil {
			return ctx.Err()
		}

		total, err := ch.startgg.client.GetPagesCount(phaseGroupId.Id, states)
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
			sets, err := ch.startgg.client.GetSets(phaseGroupId.Id, i+1, 60, states)
			if err != nil {
				log.Printf("Error getting sets: %v", err)
				continue
			}

			for _, set := range sets {
				if ctx.Err() != nil {
					break
				}

				wg.Add(1)
				go func(ctx context.Context, set startgg.Nodes, testUser sender.Participant) {
					defer wg.Done()

					if ctx.Err() != nil {
						return
					}

					// discord contact check
					p1 := ch.checkContact(set.Slots[0].Entrant.Participants)
					p2 := ch.checkContact(set.Slots[1].Entrant.Participants)

					if ctx.Err() != nil {
						return
					}

					contactP1, err := ch.searchContactDiscord(ctx, s, p1.DiscordLogin)
					if err != nil {
						log.Printf("sending message: Not finded member in discord (%v)", p1.DiscordLogin)
					} else {
						p1.DiscordID = contactP1.DiscordID
						p1.Locales = contactP1.Locales
					}

					if ctx.Err() != nil {
						return
					}

					contactP2, err := ch.searchContactDiscord(ctx, s, p2.DiscordLogin)
					if err != nil {
						log.Printf("sending message: Not finded member in discord (%v)", p2.DiscordLogin)
					} else {
						p2.DiscordID = contactP2.DiscordID
						p2.Locales = contactP2.Locales
					}

					if ctx.Err() != nil {
						return
					}

					toPlayer1 := PlayerData{
						tournament:   tournament.Name,
						setID:        set.Id,
						recipient:    contactP1, // For player1
						streamName:   set.Stream.StreamName,
						streamSourse: set.Stream.StreamSource,
						roundNum:     set.Round,
						phaseGroupId: phaseGroupId.Id,
						opponent: sender.Participant{
							DiscordID:    contactP2.DiscordID, // To player2
							DiscordLogin: p2.DiscordLogin,
							GameNickname: set.Slots[1].Entrant.Participants[0].GamerTag,
							GameID:       p2.GameID,
						},
					}
					toPlayer2 := PlayerData{
						tournament:   tournament.Name,
						setID:        set.Id,
						recipient:    contactP2, // For player2
						streamName:   set.Stream.StreamName,
						streamSourse: set.Stream.StreamSource,
						roundNum:     set.Round,
						phaseGroupId: phaseGroupId.Id,
						opponent: sender.Participant{
							DiscordID:    contactP1.DiscordID, // To player1
							DiscordLogin: p1.DiscordLogin,
							GameNickname: set.Slots[0].Entrant.Participants[0].GamerTag,
							GameID:       p1.GameID,
						},
					}

					if ch.debugMode {
						toPlayer1.recipient = testUser
						toPlayer2.recipient = testUser
					}

					if ch.debugMode || p1.DiscordLogin != "N/D" {
						ch.sendMsg(ctx, s, toPlayer1)

						select {
						case <-ctx.Done():
							return
						case <-time.After(1 * time.Second):
						}
					}

					if ch.debugMode || p2.DiscordLogin != "N/D" {
						ch.sendMsg(ctx, s, toPlayer2)

						select {
						case <-ctx.Done():
							return
						case <-time.After(1 * time.Second):
						}
					}
				}(ctx, set, testUser)
			}
			log.Printf("Checked phaseGroup(%v)", phaseGroupId)
		}
	}

	wg.Wait()
	return nil
}
