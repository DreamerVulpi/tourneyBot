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
	user         discordUser
	opponent     opponentData
}

type opponentData struct {
	discordID string
	nickname  string
	tekkenID  string
}

type discordUser struct {
	discordID string
	locales   []string
}

func (c *commandHandler) searchContactDiscord(s *discordgo.Session, nickname string) (discordUser, error) {
	member, err := s.GuildMembersSearch(c.guildID, nickname, 1)
	if err != nil {
		return discordUser{}, err
	}

	if len(member) != 1 {
		return discordUser{}, fmt.Errorf("searchContactDiscord: not finded %v", nickname)
	}

	// Get list rolesId including in locale (en is default)
	roles := []string{}
	for _, roleId := range (*member[0]).Roles {
		if roleId == c.rolesIdList.Ru {
			roles = append(roles, roleId)
		}
	}
	return discordUser{
		discordID: (*member[0]).User.ID,
		locales:   roles,
	}, nil
}

func (c *commandHandler) templateMessage(fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: c.logo,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.logoTournament,
		},
		Fields: fields,
	}
}

func (c *commandHandler) sendMessage(s *discordgo.Session, player PlayerData) {
	channel, err := s.UserChannelCreate(player.user.discordID)
	if err != nil {
		log.Println("error creating channel:", err)
		s.ChannelMessageSend(
			c.m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	link := fmt.Sprint("https://www.start.gg/", c.slug, "/set/", player.setID)

	// Check avaliable ru locale role in slice
	if len(player.user.locales) != 0 {
		c.msgInvite(s, player, channel, link, player.user.locales[0])
	} else {
		c.msgInvite(s, player, channel, link, "default")
	}
}

func (c *commandHandler) Process(s *discordgo.Session) {
	for {
		select {
		default:
			log.Println("sending messages: STARTED!")
			if err := c.SendingMessages(s); err != nil {
				return
			}
			log.Println("sending messages: DONE!")

			log.Println("sending messages: 5 minutes waiting...")
			time.Sleep(5 * time.Minute)
		case <-c.stop:
			log.Println("sending messages: STOPPED!")
			return
		}
	}
}

func (c *commandHandler) SendingMessages(s *discordgo.Session) error {
	tournament, err := c.client.GetTournament(strings.Replace(strings.SplitAfter(c.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	// Test: Change to IsDone
	if tournament.State == startgg.InProcess {
		phaseGroups, err := c.client.GetListGroups(c.slug)
		if err != nil {
			return err
		}
		for _, phaseGroup := range phaseGroups {
			state, err := c.client.GetPhaseGroupState(phaseGroup.Id)
			if err != nil {
				return err
			}
			// Test: Change to IsDone
			if state == startgg.InProcess {
				total, err := c.client.GetPagesCount(phaseGroup.Id)
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

				for i := 0; i < pages; i++ {
					sets, err := c.client.GetSets(phaseGroup.Id, pages, 60)
					if err != nil {
						log.Println(errors.New("error get sets"))
					}
					for _, set := range sets {
						if len(set.Slots) != 2 {
							continue
						}
						go func() {
							player1, err := c.searchContactDiscord(s, set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
							if err != nil {
								log.Printf("sending message: Not finded member in discord (%v)", set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
							}

							player2, err := c.searchContactDiscord(s, set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
							if err != nil {
								log.Printf("sending message: Not finded member in discord (%v)", set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
							}

							// dv, _ := c.searchContactDiscord(s, "DreamerVulpi")
							// fcuk, _ := c.searchContactDiscord(s, "fcuk_limit")

							toPlayer1 := PlayerData{
								tournament:   tournament.Name,
								setID:        set.Id,
								user:         player1, // Set player1
								streamName:   set.Stream.StreamName,
								streamSourse: set.Stream.StreamSource,
								opponent: opponentData{
									discordID: player2.discordID, // Set player2
									nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
									tekkenID:  set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
								},
							}
							toPlayer2 := PlayerData{
								tournament:   tournament.Name,
								setID:        set.Id,
								user:         player2, // Set player2
								streamName:   set.Stream.StreamName,
								streamSourse: set.Stream.StreamSource,
								opponent: opponentData{
									discordID: player1.discordID, // Set player1
									nickname:  set.Slots[0].Entrant.Participants[0].GamerTag,
									tekkenID:  set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
								},
							}

							c.sendMessage(s, toPlayer1)
							c.sendMessage(s, toPlayer2)

							log.Printf("%v vs %v -> sended! #%v", set.Slots[0].Entrant.Participants[0].GamerTag, set.Slots[1].Entrant.Participants[0].GamerTag, set.Id)
						}()
					}
					log.Printf("Checked phaseGroup(%v)", phaseGroup.Id)
				}
			}
		}
	}
	return err
}
