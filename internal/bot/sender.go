package bot

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type PlayerData struct {
	tournament   string
	setID        int64
	streamName   string
	streamSourse string
	roundNum     int
	phaseGroupId int64
	user         discordUser
	opponent     opponentData
}

type opponentData struct {
	discordID string
	nickname  string
	gameID    string
}

type discordUser struct {
	discordID string
	locales   []string
}

func (ch *commandHandler) searchContactDiscord(s *discordgo.Session, nickname string) (discordUser, error) {
	name := strings.SplitN(nickname, "#", -1)

	member, err := s.GuildMembersSearch(ch.cfg.guildID, name[0], 1)
	if err != nil {
		return discordUser{}, err
	}

	if len(member) != 1 {
		return discordUser{}, fmt.Errorf("searchContactDiscord: not finded %v", name[0])
	}

	// Get list rolesId including in locale (en is default)
	roles := []string{}
	for _, roleId := range (*member[0]).Roles {
		if roleId == ch.cfg.rolesIdList.Ru {
			roles = append(roles, roleId)
		}
	}
	return discordUser{
		discordID: strings.SplitN((*member[0]).User.ID, "#", -1)[0],
		locales:   roles,
	}, nil
}

func (ch *commandHandler) checkContact(participants []startgg.Participants) contactData {
	var discord string
	if participants == nil {
		discord = "N/D"
	} else {
		if participants[0].User.Authorizations == nil {
			value, ok := ch.discord.contacts[participants[0].GamerTag]
			if ok {
				discord = value.DiscordLogin
			} else {
				discord = "N/D"
			}
		} else {
			discord = participants[0].User.Authorizations[0].Discord
		}
	}

	var gameID string
	if participants == nil {
		gameID = "N/D"
	} else {
		if ch.cfg.tournament.Game.Name == "tekken" {
			if participants[0].ConnectedAccounts.Tekken.TekkenID == "" {
				value, ok := ch.discord.contacts[participants[0].GamerTag]
				if ok {
					gameID = value.GameID
				} else {
					gameID = "N/D"
				}
			} else {
				gameID = participants[0].ConnectedAccounts.Tekken.TekkenID
			}
		}
		if ch.cfg.tournament.Game.Name == "sf6" {
			if participants[0].ConnectedAccounts.SF6.GameID == "" {
				value, ok := ch.discord.contacts[participants[0].GamerTag]
				if ok {
					gameID = value.GameID
				} else {
					gameID = "N/D"
				}
			} else {
				gameID = participants[0].ConnectedAccounts.SF6.GameID
			}
		}
	}

	return contactData{
		DiscordLogin: discord,
		GameID:       gameID,
	}
}

func (ch *commandHandler) sendMsg(s *discordgo.Session, player PlayerData) {
	channel, err := s.UserChannelCreate(player.user.discordID)
	if err != nil {
		log.Println("error creating channel:", err)
		if _, err := s.ChannelMessageSend(
			ch.discord.msgCreate.ChannelID,
			"Something went wrong while sending the DM!",
		); err != nil {
			log.Println(err.Error())
		}
		return
	}

	link := fmt.Sprint("https://www.start.gg/", ch.slug, "/set/", player.setID)

	// Check avaliable ru locale role in slice
	if len(player.user.locales) != 0 {
		ch.msgInvite(s, player, channel, link, player.user.locales[0])
	} else {
		ch.msgInvite(s, player, channel, link, "default")
	}
}

func (ch *commandHandler) Process(s *discordgo.Session) {
	for {
		select {
		default:
			log.Println("sending messages: STARTED!")
			if err := ch.SendingMessages(s); err != nil {
				return
			}
			log.Println("sending messages: DONE!")

			log.Println("sending messages: 5 minutes waiting...")
			time.Sleep(5 * time.Minute)
		case <-ch.stopSignal:
			log.Println("sending messages: STOPPED!")
			return
		}
	}
}

func (ch *commandHandler) checkPhaseGroup(phaseGroupId int64, sets []startgg.Nodes) error {
	var min, max, minIndex, maxIndex int

	for index, set := range sets {
		if min > set.Round {
			min = set.Round
			minIndex = index
		}
		if max < set.Round {
			max = set.Round
			maxIndex = index
		}
	}

	if sets[maxIndex].FullRoundText == "Grand Final" && sets[minIndex].FullRoundText == "Losers Final" {
		log.Printf("Finded final bracket! -> %v , %v\n", min, max)
		ch.startgg.minRoundNumA = min
		ch.startgg.maxRoundNumB = max
		ch.startgg.minRoundNumB = ch.startgg.minRoundNumA + 2
		ch.startgg.maxRoundNumA = ch.startgg.maxRoundNumB - 3
		ch.startgg.finalBracketId = phaseGroupId
		return nil
	} else {
		return errors.New("not final bracket")
	}
}

func (ch *commandHandler) SendingMessages(s *discordgo.Session) error {
	tournament, err := ch.startgg.client.GetTournament(strings.Replace(strings.SplitAfter(ch.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	phaseGroups, err := ch.startgg.client.GetListGroups(ch.slug)
	if err != nil {
		return err
	}

	for _, phaseGroup := range phaseGroups {
		state, err := ch.startgg.client.GetPhaseGroupState(phaseGroup.Id)
		if err != nil {
			return err
		}
		total, err := ch.startgg.client.GetPagesCount(phaseGroup.Id)
		if err != nil {
			return err
		}
		if total == 0 {
			continue
		}

		var pages int
		if total <= 60 {
			pages = 1
		} else {
			pages = int(math.Round(float64(total / 60)))
		}

		if state == startgg.IsDone {
			for i := 0; i < pages; i++ {
				sets, err := ch.startgg.client.GetSets(phaseGroup.Id, pages, 60)
				if err != nil {
					log.Println(errors.New("error get sets"))
				}

				if err := ch.checkPhaseGroup(phaseGroup.Id, sets); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}

	for _, phaseGroup := range phaseGroups {
		state, err := ch.startgg.client.GetPhaseGroupState(phaseGroup.Id)
		if err != nil {
			return err
		}
		total, err := ch.startgg.client.GetPagesCount(phaseGroup.Id)
		if err != nil {
			return err
		}
		if total == 0 {
			continue
		}

		var pages int
		if total <= 60 {
			pages = 1
		} else {
			pages = int(math.Round(float64(total / 60)))
		}

		// Test: Set state to IsDone
		if state == startgg.InProcess {
			for i := 0; i < pages; i++ {
				sets, err := ch.startgg.client.GetSets(phaseGroup.Id, pages, 60)
				if err != nil {
					log.Println(errors.New("error get sets"))
				}

				for _, set := range sets {
					// сhecking the presence of a player in the slot
					if len(set.Slots) != 2 || len(set.Slots) == 0 {
						continue
					}
					// skip slots with empty iD
					if set.Slots[0].Entrant.Id == 0 || set.Slots[1].Entrant.Id == 0 {
						continue
					}

					go func() {
						// discord contact check
						dataPlayer1 := ch.checkContact(set.Slots[0].Entrant.Participants)
						dataPlayer2 := ch.checkContact(set.Slots[1].Entrant.Participants)

						player1, err := ch.searchContactDiscord(s, dataPlayer1.DiscordLogin)
						if err != nil {
							log.Printf("sending message: Not finded member in discord (%v)", dataPlayer1.DiscordLogin)
						}

						player2, err := ch.searchContactDiscord(s, dataPlayer2.DiscordLogin)
						if err != nil {
							log.Printf("sending message: Not finded member in discord (%v)", dataPlayer2.DiscordLogin)
						}

						// Test
						// dv, _ := ch.searchContactDiscord(s, "DreamerVulpi")

						toPlayer1 := PlayerData{
							tournament: tournament.Name,
							setID:      set.Id,
							// user: discordUser{
							// 	discordID: set.Slots[0].Entrant.Participants[0].GamerTag,
							// },
							// user: dv,
							user:         player1, // Set player1
							streamName:   set.Stream.StreamName,
							streamSourse: set.Stream.StreamSource,
							roundNum:     set.Round,
							phaseGroupId: phaseGroup.Id,
							opponent: opponentData{
								// discordID: set.Slots[1].Entrant.Participants[0].GamerTag,
								discordID: player2.discordID, // Set player2
								nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
								gameID:    dataPlayer2.GameID,
							},
						}
						toPlayer2 := PlayerData{
							tournament: tournament.Name,
							setID:      set.Id,
							// user: discordUser{
							// 	discordID: set.Slots[0].Entrant.Participants[0].GamerTag,
							// },
							// user: dv,
							user:         player2, // Set player2
							streamName:   set.Stream.StreamName,
							streamSourse: set.Stream.StreamSource,
							roundNum:     set.Round,
							phaseGroupId: phaseGroup.Id,
							opponent: opponentData{
								// discordID: set.Slots[0].Entrant.Participants[0].GamerTag,
								discordID: player1.discordID, // Set player1
								nickname:  set.Slots[0].Entrant.Participants[0].GamerTag,
								gameID:    dataPlayer1.GameID,
							},
						}

						// log.Println(toPlayer1)
						// log.Println(toPlayer2)

						if dataPlayer1.DiscordLogin != "N/D" {
							ch.sendMsg(s, toPlayer1)
							log.Printf("%v -> sended! #%v", set.Slots[0].Entrant.Participants[0].GamerTag, set.Id)
						}

						if dataPlayer2.DiscordLogin != "N/D" {
							ch.sendMsg(s, toPlayer2)
							log.Printf("%v -> sended! #%v", set.Slots[1].Entrant.Participants[0].GamerTag, set.Id)
						}
					}()
				}
				log.Printf("Checked phaseGroup(%v)", phaseGroup.Id)
			}
		}

	}
	return err
}
