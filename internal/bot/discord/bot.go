package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/config"
	"github.com/dreamervulpi/tourneyBot/internal/auth"
	"github.com/dreamervulpi/tourneyBot/internal/db/repo"
	"github.com/dreamervulpi/tourneyBot/internal/db/usecase"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
	senderUC "github.com/dreamervulpi/tourneyBot/internal/sender/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type params struct {
	guildID        string
	appID          string
	logo           string
	tournament     config.ConfigTournament
	rulesMatches   config.RulesMatches
	streamLobby    config.StreamLobby
	rolesIdList    config.ConfigRolesIdDiscord
	debugChannelID string
	debugUser      sender.Participant
}

type preparedContacts struct {
	contacts      map[string]sender.Participant
	embedContacts []*discordgo.MessageEmbed
	tourneyRole   *discordgo.Role
}

type discordHandler struct {
	auth               *auth.AuthClient
	ns                 sender.NotificationSystem
	tournamentPlatform string
	contacts           preparedContacts
	cancel             context.CancelFunc
	mutex              sync.Mutex
	cfg                params
	slug               string
	debugMode          bool
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
				MessenagerName:  dh.ns.Messenger.GetPlatformMessenagerName(),
				GameID:          dc.GameID,
				GameNickname:    dc.GameNickname,
			}

			usr, err := dh.ns.Messenger.FindContactOfParticipant(ctx, contact)
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

			sliceMessages := []*discordgo.MessageEmbed{}
			fields := []*discordgo.MessageEmbedField{}

			for nickname, dc := range dh.contacts.contacts {
				usr, err := dh.ns.Messenger.FindContactOfParticipant(ctx, dc)
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

func Init(dsAuth *auth.AuthClient, cfg config.Config, tournament config.ConfigTournament) (string, params) {
	return dsAuth.Config.ClientID, params{
		guildID:    cfg.Discord.GuildID,
		appID:      dsAuth.Config.ClientID,
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
}

func Start(dsAuth *auth.AuthClient, tourneyAuth *auth.AuthClient, conn *pgxpool.Pool, cfg config.Config, tournament config.ConfigTournament) error {
	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	appID, configTournament := Init(dsAuth, cfg, tournament)
	ctx := context.Background()

	// REFACTOR: Discord must don't know any tournament platform
	user, err := tourneyAuth.GetStartGGMe(ctx)
	if err != nil {
		return err
	}
	if len(user.ID) <= 0 {
		return fmt.Errorf("Failed get ID")
	}
	// user, err := tourneyAuth.GetChallongeMe(ctx)
	// if err != nil {
	// 	return err
	// }
	// if len(user.ID) <= 0 {
	// 	return fmt.Errorf("Failed get ID")
	// }

	log.Println(user.ID)

	ds := &DiscordSender{
		session:       session,
		participantUC: usecase.Participant{Repo: &repo.Participants{Conn: conn}},
		cfg:           configTournament,
		adminID:       user.ID,
		debugMode:     cfg.DebugMode.Mode,
	}

	ns := sender.NotificationSystem{
		Messenger:     ds,
		ParticipantUC: usecase.Participant{Repo: &repo.Participants{Conn: conn}},
		SentSetUC:     usecase.SentSet{Repo: &repo.SentSet{Conn: conn}},
		DebugMode:     cfg.DebugMode.Mode,
	}

	if cfg.DebugMode.Mode {
		me, err := dsAuth.GetDiscordMe(ctx)
		if err != nil {
			log.Printf("InitBot | Failed to get debug user: %v", err)
		} else {
			ns.TestContact = sender.Participant{
				MessenagerID:    me.ID,
				MessenagerLogin: me.Username,
				Locales:         []string{"ru"},
			}
			log.Printf("InitBot | Debug mode ON. Test contact set to: %s", me.Username)
		}
	}

	cmdHandler := discordHandler{
		auth:               dsAuth,
		cfg:                configTournament,
		ns:                 ns,
		tournamentPlatform: tournament.Platform.Platform,
		debugMode:          cfg.DebugMode.Mode,
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
	discordContacts, err := senderUC.LoadCSV(tournament.Csv.NameFile)
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
