package bot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type contactData struct {
	DiscordID    string
	DiscordLogin string
	GameID       string
}

type params struct {
	guildID      string
	appID        string
	logo         string
	tournament   config.ConfigTournament
	rulesMatches config.RulesMatches
	streamLobby  config.StreamLobby
	rolesIdList  config.ConfigRolesIdDiscord
}

// Get discord contacts from CSV file
func loadCSV(nameFile string) (map[string]contactData, error) {
	contacts := map[string]contactData{}
	f, err := os.Open("config/" + nameFile)
	if err != nil {
		log.Println(err)
		return map[string]contactData{}, err
	} else {
		if len(nameFile) != 0 {
			defer f.Close()

			csvReader := csv.NewReader(f)
			records, _ := csvReader.ReadAll()

			// Search index for get data
			var indexDiscordColumn int
			var indexGamerTagColumn int
			var indexConnectColumn int
			for index, column := range records[0] {
				parts := strings.SplitN(column, " ", -1)
				for _, part := range parts {
					if part == "Discord!" {
						indexDiscordColumn = index
					}
				}
				if column == "Short GamerTag" {
					indexGamerTagColumn = index
				}
				if column == "Connect" {
					indexConnectColumn = index
				}
			}

			for i, attendee := range records {
				if i == 0 {
					continue
				}

				var discordID string
				if len(attendee[indexDiscordColumn]) != 0 {
					discordID = attendee[indexDiscordColumn]
				} else {
					discordID = "N/D"
				}

				var gameID string
				if len(attendee[indexConnectColumn]) != 0 {
					rawGameID := strings.SplitN(attendee[indexConnectColumn], " ", -1)
					gameID = strings.ReplaceAll(rawGameID[1], ",", "")
				} else {
					gameID = "N/D"
				}

				contacts[attendee[indexGamerTagColumn]] = contactData{
					DiscordLogin: discordID,
					GameID:       gameID,
				}
			}
		}
	}

	return contacts, nil
}

func (ch *commandHandler) getDiscordContacts(s *discordgo.Session) {
	sliceMessages := []*discordgo.MessageEmbed{}
	fields := []*discordgo.MessageEmbedField{}
	counter := 0
	for nickname, dc := range ch.discord.contacts {
		if counter < 25 {
			usr, err := ch.searchContactDiscord(s, nickname)
			time.Sleep(1 * time.Second)
			if err != nil {
				log.Printf("viewContacts: %v", err.Error())
				fields = append(fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("%v", nickname), Value: fmt.Sprintf("__Discord:__\n```%v```__GameID:__\n```%v```", dc.DiscordLogin, dc.GameID), Inline: false,
				})
			} else {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("%v", nickname), Value: fmt.Sprintf("__Discord:__\n<@%v>\n__GameID:__\n```%v```", usr.discordID, dc.GameID), Inline: false,
				})
			}
			counter++
		} else {
			embed := ch.msgEmbed("", fields)
			sliceMessages = append(sliceMessages, embed)
			fields = []*discordgo.MessageEmbedField{}
			counter = 0
		}
	}

	embed := ch.msgEmbed("", fields)
	sliceMessages = append(sliceMessages, embed)
	ch.discord.embedContacts = sliceMessages
}

// Search users in server (Guild) discord from CSV file
func (ch *commandHandler) prepareContacts(s *discordgo.Session) error {
	contacts, err := os.ReadFile("contacts.json")
	if err != nil {
		log.Println("Prepare contacts from CSV...")
		for nickname, dc := range ch.discord.contacts {
			time.Sleep(1 * time.Second)
			usr, err := ch.searchContactDiscord(s, nickname)
			if err != nil {
				ch.discord.contacts[nickname] = contactData{
					DiscordID:    "N/D",
					DiscordLogin: dc.DiscordLogin,
					GameID:       dc.GameID,
				}
				log.Printf("can't find player: %v\n error: %v\n", nickname, err.Error())
				continue
			}
			ch.discord.contacts[nickname] = contactData{
				DiscordID:    usr.discordID,
				DiscordLogin: dc.DiscordLogin,
				GameID:       dc.GameID,
			}
			time.Sleep(1 * time.Second)
			err = s.GuildMemberRoleAdd(ch.cfg.guildID, usr.discordID, ch.discord.tourneyRole.ID)
			if err != nil {
				log.Println(err.Error())
			}
		}

		file, err := json.MarshalIndent(ch.discord.contacts, "", " ")
		if err != nil {
			log.Println(err.Error())
		}

		err = os.WriteFile("contacts.json", file, 0644)
		if err != nil {
			log.Println(err.Error())
		}

		log.Println("Done!")
	} else {
		err := json.Unmarshal(contacts, &ch.discord.contacts)
		if err != nil {
			return err
		}

		log.Println("Loaded contact.json file")
	}

	contactsEmbed, err := os.ReadFile("contactsEmbed.json")
	if err != nil {
		if len(ch.discord.contacts) != 0 {
			log.Println("Generate contact.json file...")

			ch.getDiscordContacts(s)

			file, err := json.MarshalIndent(ch.discord.embedContacts, "", " ")
			if err != nil {
				return err
			}

			err = os.WriteFile("contactsEmbed.json", file, 0644)
			if err != nil {
				return err
			}
		} else {
			log.Println("Error: List discord contacts is empty")
		}
	} else {
		err := json.Unmarshal(contactsEmbed, &ch.discord.embedContacts)
		if err != nil {
			return err
		}
		log.Println("Loaded contactEmbed.json file")
	}
	return nil
}

func (ch *commandHandler) createTourneyRole(session *discordgo.Session) error {
	rolesServer, err := session.GuildRoles(ch.cfg.guildID)
	if err != nil {
		return err
	}

	var checker bool

	// check available role in guild (server) discord
	for _, r := range rolesServer {
		if r.Name == "Tourney Role" {
			checker = true
			ch.discord.tourneyRole = r
			log.Println("Finded role in server! Saved to program.")
		}
	}

	if !checker {
		color := 16711680
		hoist := true
		mentionable := true
		var prms int64 = 0x0000000000000800 | 0x0000000000000400

		rslt, err := session.GuildRoleCreate(ch.cfg.guildID, &discordgo.RoleParams{
			Name:        "Tourney Role",
			Color:       &color,
			Hoist:       &hoist,
			Mentionable: &mentionable,
			Permissions: &prms,
		})

		if err != nil {
			log.Println(err.Error())
		}

		ch.discord.tourneyRole = rslt

		log.Println("Tourney role successfuly created in server!")
	}

	return nil
}

func (ch *commandHandler) deleteTourneyRole(session *discordgo.Session) error {
	rolesServer, err := session.GuildRoles(ch.cfg.guildID)
	if err != nil {
		return err
	}

	// check available role in guild (server) discord
	for _, r := range rolesServer {
		if r.Name == "Tourney Role" {
			err := session.GuildRoleDelete(ch.cfg.guildID, r.ID)
			if err != nil {
				return err
			}
			log.Println("Tourney role successfuly deleted from server!")
			break
		}
	}

	return nil
}

func Start(cfg config.Config, tournament config.ConfigTournament) error {
	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	client := startgg.NewClient(cfg.Startgg.Token, &http.Client{
		Timeout: time.Second * 10,
	})

	cfgCmdHadnler := params{
		guildID:    cfg.Discord.GuildID,
		appID:      cfg.Discord.AppID,
		tournament: tournament,
		rulesMatches: config.RulesMatches{
			StandardFormat: tournament.Rules.StandardFormat,
			FinalsFormat:   tournament.Rules.FinalsFormat,
			Rounds:         tournament.Rules.Rounds,
			Duration:       tournament.Rules.Duration,
			Crossplatform:  tournament.Rules.Crossplatform,
			Stage:          tournament.Rules.Stage,
		},
		streamLobby: config.StreamLobby{
			Area:          tournament.Stream.Area,
			Language:      tournament.Stream.Language,
			Crossplatform: tournament.Stream.Crossplatform,
			Conn:          tournament.Stream.Conn,
			Passcode:      tournament.Stream.Passcode,
		},
		rolesIdList: cfg.Roles,
		logo:        "https://i.imgur.com/n9SG5IL.png",
	}

	cmdHandler := commandHandler{
		stopSignal: make(chan struct{}),
		startgg: strtgg{
			client: client,
		},
		cfg: cfgCmdHadnler,
	}

	commandHandlers["check"] = cmdHandler.viewData
	commandHandlers["start-sending"] = cmdHandler.startSending
	commandHandlers["stop-sending"] = cmdHandler.stopSending
	commandHandlers["set-event"] = cmdHandler.setEvent
	commandHandlers["edit-rules"] = cmdHandler.editRuleMatches
	commandHandlers["edit-stream-lobby"] = cmdHandler.editStreamLobby
	commandHandlers["edit-logo-tournament"] = cmdHandler.editLogoTournament

	var trigger bool
	discordContacts, err := loadCSV(tournament.Csv.NameFile)
	cmdHandler.discord.contacts = discordContacts
	if err != nil {
		log.Println("CSV file isn't loaded. Commands: contacts and roles unavailable. Autofill empty data unavailable.")
		trigger = true
	} else {

		err = cmdHandler.createTourneyRole(session)
		if err != nil {
			return err
		}
		err = cmdHandler.prepareContacts(session)
		if err != nil {
			return err
		}
		commandHandlers["contacts"] = cmdHandler.viewContacts
		commandHandlers["roles"] = cmdHandler.roles
	}

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	log.Println("adding commands...")
	commands := cmdHandler.commands()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		if command.Name == "roles" && trigger || command.Name == "contacts" && trigger {
			continue
		}
		cmd, err := session.ApplicationCommandCreate(cfg.Discord.AppID, cfg.Discord.GuildID, command)
		log.Printf("%v\n", command.Name)
		if err != nil {
			log.Printf("can't create '%v' command: %v\n", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	log.Println("the bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press Ctrl+C to exit")
	<-stop

	log.Println("removing commands...")
	if err := cmdHandler.deleteTourneyRole(session); err != nil {
		return err
	}

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(cfg.Discord.AppID, cfg.Discord.GuildID, v.ID)
		log.Printf("%v\n", v.Name)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
	log.Println("gracefully shutting down!")
	return nil
}
