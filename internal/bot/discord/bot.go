package discord

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"errors"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/db/repo"
	"github.com/dreamervulpi/tourneyBot/internal/db/usecase"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	"github.com/dreamervulpi/tourneyBot/startgg"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type sender.Participant struct {
// 	MessenagerID    string
// 	MessenagerLogin string
// 	GameID          string
// }

type params struct {
	guildID        string
	appID          string
	logo           string
	tournament     config.ConfigTournament
	rulesMatches   config.RulesMatches
	streamLobby    config.StreamLobby
	rolesIdList    config.ConfigRolesIdDiscord
	debugChannelID string
}

type preparedContacts struct {
	contacts      map[string]sender.Participant
	embedContacts []*discordgo.MessageEmbed
	tourneyRole   *discordgo.Role
}

type strtgg struct {
	client         *startgg.Client
	finalBracketId int64
	minRoundNumA   int
	minRoundNumB   int
	maxRoundNumA   int
	maxRoundNumB   int
}

// TODO: Change startgg on universal
type discordHandler struct {
	auth          *auth.AuthClient
	msgSender     sender.NotificationSender
	contacts      preparedContacts
	cancel        context.CancelFunc
	startgg       strtgg
	mutex         sync.Mutex
	sentSetUC     usecase.SentSet
	participantUC usecase.Participant
	cfg           params
	slug          string
	debugMode     bool
}

// Get discord contacts from CSV file Startgg
func loadCSV(nameFile string) (map[string]sender.Participant, error) {
	if nameFile == "" {
		return nil, errors.New("loadCSV: filename is empty")
	}

	f, err := os.Open("config/" + nameFile)
	if err != nil {
		return map[string]sender.Participant{}, fmt.Errorf("loadCSV: open file, %v", err)
	}
	defer f.Close() //nolint:errcheck

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return map[string]sender.Participant{}, fmt.Errorf("loadCSV: read CSV, %v", err)
	}

	if len(records) == 0 {
		return map[string]sender.Participant{}, nil
	}

	// Search index for get data
	var idxDiscordColumn, idxGamerTagColumn, idxConnectColumn int
	for index, column := range records[0] {
		if strings.Contains(column, "Discord!") {
			idxDiscordColumn = index
		}
		if column == "Short GamerTag" {
			idxGamerTagColumn = index
		}
		if column == "Connect" {
			idxConnectColumn = index
		}
	}

	contacts := make(map[string]sender.Participant, len(records)-1)
	for i := 1; i < len(records); i++ {
		attendee := records[i]

		discordID := "N/D"
		if val := attendee[idxDiscordColumn]; val != "" {
			discordID = val
		}

		gameID := "N/D"
		if val := attendee[idxConnectColumn]; val != "" {
			rawGameID := strings.Split(attendee[idxConnectColumn], " ")
			if len(rawGameID) >= 2 {
				gameID = strings.ReplaceAll(rawGameID[1], ",", "")
			}
		}

		gameNickname := "N/D"
		if val := attendee[idxGamerTagColumn]; val != "" {
			gameNickname = val
		}

		key := attendee[idxGamerTagColumn]
		if key != "" {
			contacts[key] = sender.Participant{
				MessenagerLogin: discordID,
				GameID:          gameID,
				GameNickname:    gameNickname,
			}
		}
	}

	return contacts, nil
}

func (dh *discordHandler) getDiscordContacts(ctx context.Context, s *discordgo.Session) {
	sliceMessages := []*discordgo.MessageEmbed{}
	fields := []*discordgo.MessageEmbedField{}

	for nickname, dc := range dh.contacts.contacts {
		usr, err := dh.searchContactDiscord(ctx, s, dc.MessenagerLogin, nickname)
		time.Sleep(1 * time.Second)
		field := &discordgo.MessageEmbedField{
			Name:   nickname,
			Inline: false,
		}

		if err != nil {
			field.Value = fmt.Sprintf("__Discord:__\n```%v```__GameID:__\n```%v```", dc.MessenagerLogin, dc.GameID)
		} else {
			field.Value = fmt.Sprintf("__Discord:__\n<@%v>\n__GameID:__\n```%v```", usr.MessenagerID, dc.GameID)
		}

		fields = append(fields, field)

		if len(fields) == 25 {
			sliceMessages = append(sliceMessages, dh.msgEmbed("", fields, ColorSystem))
			fields = []*discordgo.MessageEmbedField{}
		}
	}

	if len(fields) > 0 {
		sliceMessages = append(sliceMessages, dh.msgEmbed("", fields, ColorSystem))
	}

	dh.contacts.embedContacts = sliceMessages
}

// Search users in server (Guild) discord from CSV file
func (dh *discordHandler) prepareContacts(ctx context.Context, s *discordgo.Session) error {
	contacts, err := os.ReadFile("contacts.json")
	if err != nil {
		log.Println("Prepare contacts from CSV...")
		for nickname, dc := range dh.contacts.contacts {
			time.Sleep(1 * time.Second)
			contact := sender.Participant{
				MessenagerID:    "N/D",
				MessenagerLogin: dc.MessenagerLogin,
				GameID:          dc.GameID,
				GameNickname:    dc.GameNickname,
			}

			usr, err := dh.searchContactDiscord(ctx, s, dc.MessenagerLogin, nickname)
			if err != nil {
				dh.contacts.contacts[nickname] = contact
				log.Printf("can't find player: %v\n error: %v\n", nickname, err.Error())
				continue
			}
			contact.MessenagerID = usr.MessenagerID
			dh.contacts.contacts[nickname] = contact

			time.Sleep(1 * time.Second)
			if usr.MessenagerID != "000000000000000000" && usr.MessenagerID != "N/D" {
				err = s.GuildMemberRoleAdd(dh.cfg.guildID, usr.MessenagerID, dh.contacts.tourneyRole.ID)
				if err != nil {
					log.Printf("prepareContacts | discord API Error (RoleAdd) for %v: %v", nickname, err)
					continue
				}
			} else {
				log.Printf("prepareContacts | skip roleAdd for %v: user not on server (Mock ID used)\n", nickname)
			}
		}

		file, err := json.MarshalIndent(dh.contacts.contacts, "", " ")
		if err != nil {
			log.Println(err.Error())
		}

		err = os.WriteFile("contacts.json", file, 0644)
		if err != nil {
			log.Println(err.Error())
		}

		log.Println("Done!")
	} else {
		err := json.Unmarshal(contacts, &dh.contacts.contacts)
		if err != nil {
			return err
		}

		log.Println("Loaded contact.json file")
	}

	contactsEmbed, err := os.ReadFile("contactsEmbed.json")
	if err != nil {
		if len(dh.contacts.contacts) != 0 {
			log.Println("Generate contact.json file...")

			dh.getDiscordContacts(ctx, s)

			file, err := json.MarshalIndent(dh.contacts.embedContacts, "", " ")
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
		err := json.Unmarshal(contactsEmbed, &dh.contacts.embedContacts)
		if err != nil {
			return err
		}
		log.Println("Loaded contactEmbed.json file")
	}
	return nil
}

func (s *discordHandler) createTourneyRole(session *discordgo.Session) error {
	rolesServer, err := session.GuildRoles(s.cfg.guildID)
	if err != nil {
		return err
	}

	var checker bool

	// check available role in guild (server) discord
	for _, r := range rolesServer {
		if r.Name == "Tourney Role" {
			checker = true
			s.contacts.tourneyRole = r
			log.Println("createTourneyRole | Finded role in server! Saved to program")
		}
	}

	if !checker {
		color := 16711680
		hoist := true
		mentionable := true
		var prms int64 = 0x0000000000000800 | 0x0000000000000400

		rslt, err := session.GuildRoleCreate(s.cfg.guildID, &discordgo.RoleParams{
			Name:        "Tourney Role",
			Color:       &color,
			Hoist:       &hoist,
			Mentionable: &mentionable,
			Permissions: &prms,
		})

		if err != nil {
			log.Println(err.Error())
		}

		s.contacts.tourneyRole = rslt

		log.Println("Tourney role successfuly created in server!")
	}

	return nil
}

func (s *discordHandler) deleteTourneyRole(session *discordgo.Session) error {
	rolesServer, err := session.GuildRoles(s.cfg.guildID)
	if err != nil {
		return err
	}

	// check available role in guild (server) discord
	for _, r := range rolesServer {
		if r.Name == "Tourney Role" {
			err := session.GuildRoleDelete(s.cfg.guildID, r.ID)
			if err != nil {
				return err
			}
			log.Println("Tourney role successfuly deleted from server!")
			break
		}
	}

	return nil
}

func Start(stClient *http.Client, dsAuth *auth.AuthClient, pool *pgxpool.Pool, cfg config.Config, tournament config.ConfigTournament) error {
	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	appID := dsAuth.Config.ClientID

	configTournament := params{
		guildID:    cfg.Discord.GuildID,
		appID:      appID,
		tournament: tournament,
		rulesMatches: config.RulesMatches{
			StandardFormat: tournament.Rules.StandardFormat,
			FinalsFormat:   tournament.Rules.FinalsFormat,
			Rounds:         tournament.Rules.Rounds,
			Duration:       tournament.Rules.Duration,
			Crossplatform:  tournament.Rules.Crossplatform,
			Stage:          tournament.Rules.Stage,
			Waiting:        tournament.Rules.Waiting,
		},
		streamLobby: config.StreamLobby{
			Area:          tournament.Stream.Area,
			Language:      tournament.Stream.Language,
			Crossplatform: tournament.Stream.Crossplatform,
			Conn:          tournament.Stream.Conn,
			Passcode:      tournament.Stream.Passcode,
		},
		rolesIdList:    cfg.Roles,
		logo:           "https://i.imgur.com/n9SG5IL.png",
		debugChannelID: cfg.Discord.DebugChannelID,
	}

	sentSetUsecase := usecase.SentSet{Repo: &repo.SentSet{Conn: pool}}
	participantsUsecase := usecase.Participant{Repo: &repo.Participants{Conn: pool}}

	startggClient := startgg.NewClient(stClient)

	dsSender := &DiscordSender{
		session: session,
		config:  configTournament,
		// TODO: Delete
		startgg: strtgg{
			client: startggClient,
		},
		debugChannelID: cfg.Discord.DebugChannelID,
	}

	cmdHandler := discordHandler{
		startgg: strtgg{
			client: startggClient,
		},
		auth:          dsAuth,
		cfg:           configTournament,
		msgSender:     dsSender,
		sentSetUC:     sentSetUsecase,
		participantUC: participantsUsecase,
		debugMode:     cfg.DebugMode.Mode,
	}

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commandHandlers["check"] = cmdHandler.viewData
	commandHandlers["start-sending"] = cmdHandler.startSending
	commandHandlers["stop-sending"] = cmdHandler.stopSending
	commandHandlers["set-event"] = cmdHandler.setEvent
	commandHandlers["edit-rules"] = cmdHandler.editRuleMatches
	commandHandlers["edit-stream-lobby"] = cmdHandler.editStreamLobby
	commandHandlers["edit-logo-tournament"] = cmdHandler.editLogoTournament

	var trigger bool
	discordContacts, err := loadCSV(tournament.Csv.NameFile)
	cmdHandler.contacts.contacts = discordContacts
	if err != nil {
		log.Println("CSV file isn't loaded. Commands: contacts and roles unavailable. Autofill empty data unavailable.")
		trigger = true
	} else {
		err = cmdHandler.createTourneyRole(session)
		if err != nil {
			return err
		}
		err = cmdHandler.prepareContacts(context.Background(), session)
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
		cmd, err := session.ApplicationCommandCreate(appID, cfg.Discord.GuildID, command)
		log.Printf("%v\n", command.Name)
		if err != nil {
			log.Printf("can't create '%v' command: %v\n", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer session.Close() //nolint:errcheck

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
		err := session.ApplicationCommandDelete(appID, cfg.Discord.GuildID, v.ID)
		log.Printf("%v\n", v.Name)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
	log.Println("gracefully shutting down!")
	return nil
}
