package bot

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type playerData struct {
	setID     int64
	discordID string
	opponent  opponentData
}

type opponentData struct {
	discordID string
	nickname  string
	tekkenID  string
}

func (c *commandHandler) searchContactDiscord(s *discordgo.Session, nickname string) (string, error) {
	// if !server() {
	// 	return "", errors.New("server ID is empty")
	// }
	user, err := s.GuildMembersSearch(c.guildID, nickname, 1)
	if err != nil {
		return "", err
	}
	return (*user[0]).User.ID, nil
}

func (c *commandHandler) sendMessage(s *discordgo.Session, player playerData) {
	channel, err := s.UserChannelCreate(player.discordID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			c.m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	link := fmt.Sprint("https://www.start.gg/", c.slug, "/set/", player.setID)
	invite := fmt.Sprintf(c.templateInviteMessage, "турик", player.opponent.nickname, player.opponent.tekkenID, player.opponent.discordID, link)

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
			fmt.Println("Stopped.")
			return nil
		default:
			fmt.Println("Start sending messages...")
			if err := c.SendProcess(s); err != nil {
				return err
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (c *commandHandler) SendProcess(s *discordgo.Session) error {

	// phaseGroups, err := startgg.GetListPhaseGroups(Slug)
	phaseGroups, err := c.client.GetListGroups(c.slug)
	if err != nil {
		return err
	}

	// TODO: GetListPhaseGroups
	// for _, pgs := range results.Data.Event.PhaseGroups {
	// 	pgs.Id
	// }

	groupId := phaseGroups[0].Id

	state, err := c.client.GetPhaseGroupState(groupId)
	if err != nil {
		return err
	}

	if state == startgg.IsDone {
		total, err := c.client.GetPagesCount(groupId)
		if err != nil {
			return err
		}

		var pages int

		if total <= 60 {
			pages = 1
		} else {
			pages = int(math.Round(float64(total / 60)))
		}

		fmt.Println(total, "/", 60, "=", pages)

		for i := 0; i < pages; i++ {
			sets, err := c.client.GetSets(groupId, pages, 60)
			for _, set := range sets {

				// checkIn := fmt.Sprint("https://www.start.gg/", Slug, "/set/", set.Id)
				// checkIn := fmt.Sprint("https://www.start.gg/", "tournament/wild-hunters-1/event/main-online-crossplatform-event", "/set/", set.Id)
				// fmt.Println("generated CheckIn: ", checkIn)
				// fmt.Println("set ID: ", set.Id)
				// fmt.Println("state set: ", set.State)

				// player1Discord, _ := searchContactDiscord(s, set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
				// player2Discord, _ := searchContactDiscord(s, set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)

				dv, _ := c.searchContactDiscord(s, "DreamerVulpi")
				fcuk, _ := c.searchContactDiscord(s, "fcuk_limit")

				toPlayer1 := playerData{
					setID:     set.Id,
					discordID: dv,
					opponent: opponentData{
						discordID: fcuk,
						nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
						tekkenID:  set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
					},
				}

				// log.Printf("player 1 | Discord: ", player1Discord)

				c.sendMessage(s, toPlayer1)

				toPlayer2 := playerData{
					setID:     set.Id,
					discordID: dv,
					opponent: opponentData{
						discordID: fcuk,
						nickname:  set.Slots[0].Entrant.Participants[0].GamerTag,
						tekkenID:  set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
					},
				}

				// log.Printf("player 2 | Discord: ", player2Discord)

				c.sendMessage(s, toPlayer2)

				fmt.Println("sended messages..")
				// fmt.Println(player1.User.ID)
				// fmt.Println("player 1 | ID: ", set.Slots[0].Entrant.Id)
				// fmt.Println("player 1 | Nickname: ", set.Slots[0].Entrant.Participants[0].GamerTag)
				// fmt.Println("player 1 | TEKKEN ID ", set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID)
				// fmt.Println("player 1 | Discord: ", set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)

				// fmt.Println("player 2 | ID: ", set.Slots[1].Entrant.Id)
				// fmt.Println("player 2 | Nickname: ", set.Slots[1].Entrant.Participants[0].GamerTag)
				// fmt.Println("player 2 | TEKKEN ID ", set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID)
				// fmt.Println("player 2 | Discord: ", set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
			}
			if err != nil {
				fmt.Println(errors.New("error get sets"))
			}
		}

	}

	return err
}
