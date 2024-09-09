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
	tournament   string
	setID        int64
	discordID    string
	streamSourse string
	opponent     opponentData
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

// TODO: Refactor code
// TODO: Add can use different languages
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
	var t discordgo.MessageEmbed
	if len(player.streamSourse) == 0 {
		var crossplay string
		if c.dataLobby.Local.CrossPlay != "true" {
			crossplay = "отключена"
		} else {
			crossplay = "включена"
		}
		var conn string
		if c.dataLobby.Local.Conn == "no restrictions" {
			conn = "нет ограничений"
		}
		var lang string
		if c.dataLobby.Local.Language == "any" {
			lang = "любой"
		}
		var area string
		if c.dataLobby.Local.Area == "any" {
			area = "любой"
		}
		t = discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Турнир **%v**", player.tournament),
			Description: "Приглашение на турнир со всей необходимой информацией.\n\n*Это сообщение сгенерировано автоматически. Отвечать на него не нужно. В случае вопросов или помощи обращайтесь к помощникам организатора.*",
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://i.imgur.com/AfFp7pu.png",
				URL:     "https://github.com/DreamerVulpi/tourneybot",
				Name:    "TourneyBot",
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "**Данные твоего оппонента**", Value: ""},
				{Name: "**Никнейм**", Value: fmt.Sprintf("```%v```", player.opponent.nickname), Inline: true},
				{Name: "**Tekken ID**", Value: fmt.Sprintf("```%v```", player.opponent.tekkenID), Inline: true},
				{Name: "**Discord**", Value: fmt.Sprintf("<@%v>", player.opponent.discordID), Inline: true},

				{Name: "**Ссылка на check-in**", Value: link},
				{Name: "У вас есть 10 минут чтобы отметиться до автоматической дисквалификации", Value: ""},

				{Name: "**Настройки лобби согласно правилам**", Value: ""},
				{Name: "", Value: ""},
				{Name: "**Регион**", Value: area, Inline: true},
				{Name: "**Язык**", Value: lang, Inline: true},
				{Name: "**Тип соединения**", Value: conn, Inline: true},
				{Name: ""},

				{Name: "**Настройки боев**", Value: ""},
				{Name: ""},
				{Name: "**Формат**", Value: fmt.Sprintf("ФТ%v", c.dataLobby.Local.Victory) + fmt.Sprintf(" (До %v побед)", c.dataLobby.Local.Victory), Inline: true},

				{Name: "**Карта**", Value: "Выбирается случайным образом ВСЕГДА если оппонент не продолжил сет", Inline: true},
				{Name: ""},

				{Name: "**Раундов в 1 матче**", Value: fmt.Sprintf("%v", c.dataLobby.Local.Rounds), Inline: true},
				{Name: "**Время в 1 раунде**", Value: fmt.Sprintf("%v", c.dataLobby.Local.Duration) + " секунд", Inline: true},
				{Name: "**Кроссплатформенная игра**", Value: crossplay, Inline: true},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://i.imgur.com/AfFp7pu.png",
			},
		}
	}
	if len(player.streamSourse) > 0 {
		var lang string
		if c.dataLobby.Local.Language == "any" {
			lang = "любой"
		}
		var area string
		if c.dataLobby.Local.Area == "any" {
			area = "любой"
		}
		var crossplay string
		if c.dataLobby.Local.CrossPlay != "true" {
			crossplay = "отключена"
		} else {
			crossplay = "включена"
		}
		var conn string
		if c.dataLobby.Local.Conn == "no restrictions" {
			conn = "нет ограничений"
		}
		t = discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Турнир: **%v**", player.tournament),
			Description: "Приглашение на матч проводящийся на прямой трансляции. Необходимо зайти в ниже указанное лобби и ожидать команды организатора на стриме дальнейшних действий.\n\n*Это сообщение сгенерировано автоматически. Отвечать на него не нужно. В случае вопросов или помощи обращайтесь к помощникам организатора.*",
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://i.imgur.com/AfFp7pu.png",
				URL:     "https://github.com/DreamerVulpi/tourneybot",
				Name:    "TourneyBot",
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "**Ссылка на check-in**", Value: ""},
				{Name: "*У вас есть 10 минут чтобы отметиться до автоматической дисквалификации*"},
				{Name: "", Value: link, Inline: true},

				{Name: "**Параметры для поиска лобби**", Value: ""},
				{Name: "", Value: ""},
				{Name: "**Регион**", Value: area, Inline: true},
				{Name: "**Язык**", Value: lang, Inline: true},
				{Name: "**Тип соединения**", Value: conn, Inline: true},
				{Name: "**Кроссплатформенная игра**", Value: crossplay, Inline: true},
				{Name: "**Пароль**", Value: fmt.Sprintf("```%v```", c.dataLobby.Stream.Passcode), Inline: true},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://i.imgur.com/AfFp7pu.png",
			},
		}
	}

	_, err = s.ChannelMessageSendEmbed(channel.ID, &t)
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			c.m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

func (c *commandHandler) SendingMessages(s *discordgo.Session) {
	for {
		select {
		default:
			log.Println("sending messages: STARTED!")
			if err := c.SendProcess(s); err != nil {
				return
			}
			log.Println("sending messages: DONE!")
			// TODO: TIMER 5 MINUTES
			log.Println("sending messages: 5 minutes waiting...")
			time.Sleep(2 * time.Second)
		case <-c.stop:
			log.Println("sending messages: STOPPED!")
			return
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

				for i := 0; i < pages; i++ {
					sets, err := c.client.GetSets(phaseGroup.Id, pages, 60)
					if err != nil {
						log.Println(errors.New("error get sets"))
					}
					for _, set := range sets {
						go func() {
							// TODO: player1, err := c.searchContactDiscord(s, set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
							// if err != nil {
							// 	log.Printf("sending message: Not finded member in discord (%v)", set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
							// }

							// TODO: player2, err := c.searchContactDiscord(s, set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
							// if err != nil {
							// 	log.Printf("sending message: Not finded member in discord (%v)", set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
							// }
							dv, _ := c.searchContactDiscord(s, "DreamerVulpi")
							fcuk, _ := c.searchContactDiscord(s, "fcuk_limit")

							toPlayer1 := playerData{
								tournament:   tournament.Name,
								setID:        set.Id,
								discordID:    dv, // TODO: Set player1
								streamSourse: set.Stream.StreamSource,
								opponent: opponentData{
									discordID: fcuk, // TODO: Set player2
									nickname:  set.Slots[1].Entrant.Participants[0].GamerTag,
									tekkenID:  set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID,
								},
							}
							toPlayer2 := playerData{
								tournament:   tournament.Name,
								setID:        set.Id,
								discordID:    dv, // TODO: Set player2
								streamSourse: set.Stream.StreamSource,
								opponent: opponentData{
									discordID: fcuk, // TODO: Set player1
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
