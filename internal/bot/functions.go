package bot

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/internal/startgg/functions"
)

type player struct {
	setID     int64
	discordID string
	opponent  opponent
}

type opponent struct {
	discordID string
	nickname  string
	tekkenID  string
}

var (
	SendData player
	Pages    int
)

func SetSendData(data player) {
	SendData = data
}

func SetPages(number int) {
	Pages = number
}

func searchContactDiscord(s *discordgo.Session, nickname string) (string, error) {
	if !server() {
		return "", errors.New("server ID is empty")
	}
	user, err := s.GuildMembersSearch(GuildID, nickname, 1)
	if err != nil {
		return "", err
	}
	return (*user[0]).User.ID, nil
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel, err := s.UserChannelCreate(SendData.discordID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	link := fmt.Sprint("https://www.start.gg/", "tournament/wild-hunters-1/event/main-online-crossplatform-event", "/set/", SendData.setID)
	invite := fmt.Sprintf(templateInviteMessage, "турик", SendData.opponent.nickname, SendData.opponent.tekkenID, SendData.opponent.discordID, link)

	_, err = s.ChannelMessageSend(channel.ID, invite)
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

func SendProcess(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// if !slug() {
	// 	return errors.New("slug data is empty.")
	// }

	// phaseGroups, err := startgg.GetListPhaseGroups(Slug)
	phaseGroups, err := functions.GetListGroups("tournament/wild-hunters-1/event/main-online-crossplatform-event")
	if err != nil {
		return err
	}

	// TODO: GetListPhaseGroups
	// for _, pgs := range results.Data.Event.PhaseGroups {
	// 	pgs.Id
	// }

	groupId := phaseGroups[0].Id

	state, err := functions.GetPhaseGroupState(groupId)
	if err != nil {
		return err
	}

	// State = 1 (not started) v
	// State = 2 (in process)
	// State = 3 (done)

	if state == 3 {
		// TODO: Query for pages
		total, err := functions.GetPagesCount(groupId)
		if err != nil {
			return err
		}

		if total <= 60 {
			SetPages(1)
		} else {
			SetPages(int(math.Round(float64(total / 60))))
		}

		fmt.Println(total, "/", 60, "=", Pages)

		for i := 0; i < Pages; i++ {
			sets, err := functions.GetSets(groupId, Pages, 60)
			for _, set := range sets {
				time.Sleep(1 * time.Second)
				// checkIn := fmt.Sprint("https://www.start.gg/", Slug, "/set/", set.Id)
				// checkIn := fmt.Sprint("https://www.start.gg/", "tournament/wild-hunters-1/event/main-online-crossplatform-event", "/set/", set.Id)
				// fmt.Println("generated CheckIn: ", checkIn)
				// fmt.Println("set ID: ", set.Id)
				// fmt.Println("state set: ", set.State)

				// player1Discord, _ := searchContactDiscord(s, set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
				// player2Discord, _ := searchContactDiscord(s, set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)

				dv, _ := searchContactDiscord(s, "DreamerVulpi")
				fcuk, _ := searchContactDiscord(s, "fcuk_limit")

				sendtoPlayer1 := player{
					setID:     set.Id,
					discordID: dv,
					opponent: opponent{
						discordID: fcuk,
						nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
						tekkenID:  set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
					},
				}

				// log.Printf("player 1 | Discord: ", player1Discord)

				SetSendData(sendtoPlayer1)

				sendMessage(s, m)

				sendtoPlayer2 := player{
					setID:     set.Id,
					discordID: dv,
					opponent: opponent{
						discordID: fcuk,
						nickname:  set.Slots[0].Entrant.Participants[0].GamerTag,
						tekkenID:  set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
					},
				}

				time.Sleep(1 * time.Second)

				// log.Printf("player 2 | Discord: ", player2Discord)

				SetSendData(sendtoPlayer2)
				sendMessage(s, m)

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
