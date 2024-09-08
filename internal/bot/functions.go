package bot

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type playerData struct {
	tournament string
	setID      int64
	discordID  string
	opponent   opponentData
}

type opponentData struct {
	discordID string
	nickname  string
	tekkenID  string
}

func (c *commandHandler) searchContactDiscord(s *discordgo.Session, nickname string) (string, error) {
	user, err := s.GuildMembersSearch(c.guildID, nickname, 1)
	if err != nil {
		return "", err
	}
	return (*user[0]).User.ID, nil
}

func (c *commandHandler) sendMessage(s *discordgo.Session, player playerData) {
	channel, err := s.UserChannelCreate(player.discordID)
	if err != nil {
		log.Println("error creating channel:", err)
		s.ChannelMessageSend(
			c.m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	link := fmt.Sprint("https://www.start.gg/", c.slug, "/set/", player.setID)
	invite := fmt.Sprintf(c.templateInviteMessage, player.tournament, player.opponent.nickname, player.opponent.tekkenID, player.opponent.discordID, link)

	_, err = s.ChannelMessageSend(channel.ID, invite)
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			c.m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

func (c *commandHandler) SendingMessages(s *discordgo.Session) error {
	for {
		select {
		case <-c.stop:
			log.Println("sending messages: STOPPED!")
			return nil
		default:
			log.Println("sending messages: STARTED!")
			if err := c.SendProcess(s); err != nil {
				return err
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (c *commandHandler) SendProcess(s *discordgo.Session) error {
	tournament, err := c.client.GetTournament(strings.Replace(strings.SplitAfter(c.slug, "/")[1], "/", "", 1))
	if err != nil {
		return err
	}

	// TODO: Change to InProcess
	if tournament.State == startgg.IsDone {
		phaseGroups, err := c.client.GetListGroups(c.slug)
		if err != nil {
			return err
		}
		for _, phaseGroup := range phaseGroups {
			state, err := c.client.GetPhaseGroupState(phaseGroup.Id)
			if err != nil {
				return err
			}
			// TODO: Change to InProcess
			if state == startgg.IsDone {
				total, err := c.client.GetPagesCount(phaseGroup.Id)
				if err != nil {
					return err
				}

				var pages int

				if total == 0 {
					continue
				}

				if total <= 60 {
					pages = 1
				} else {
					pages = int(math.Round(float64(total / 60)))
				}

				fmt.Println(total, "/", 60, "=", pages)

				for i := 0; i < pages; i++ {

					sets, err := c.client.GetSets(phaseGroup.Id, pages, 60)
					if err != nil {
						log.Println(errors.New("error get sets"))
					}
					for _, set := range sets {
						// TODO: Change to NotStarted
						if set.State == startgg.IsDone {
							// TODO: player1, _ := c.searchContactDiscord(s, set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
							// TODO: player2, _ := c.searchContactDiscord(s, set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
							dv, _ := c.searchContactDiscord(s, "DreamerVulpi")
							fcuk, _ := c.searchContactDiscord(s, "fcuk_limit")

							toPlayer1 := playerData{
								tournament: tournament.Name,
								setID:      set.Id,
								discordID:  dv, // TODO: Set player1
								opponent: opponentData{
									discordID: fcuk, // TODO: Set player2
									nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
									tekkenID:  set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
								},
							}

							c.sendMessage(s, toPlayer1)

							toPlayer2 := playerData{
								tournament: tournament.Name,
								setID:      set.Id,
								discordID:  dv, // TODO: Set player2
								opponent: opponentData{
									discordID: fcuk, // TODO: Set player1
									nickname:  set.Slots[0].Entrant.Participants[0].GamerTag,
									tekkenID:  set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
								},
							}

							c.sendMessage(s, toPlayer2)

							log.Printf("%v vs %v -> sended! #%v", set.Slots[0].Entrant.Participants[0].GamerTag, set.Slots[1].Entrant.Participants[0].GamerTag, set.Id)
						}
					}
				}
			}
		}
	}
	return err
}
