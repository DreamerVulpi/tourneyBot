package bot

import (
	"encoding/json"
	"fmt"
	"log"

	"net/url"
	"strings"

	"regexp"

	"os"

	"github.com/bwmarrin/discordgo"
)

func (ch *commandHandler) controlRole(s *discordgo.Session, arg string) []*discordgo.MessageEmbed {
	var embed []*discordgo.MessageEmbed
	if len(ch.discord.contacts) != 0 {
		if arg == "give" {
			for _, usr := range ch.discord.contacts {
				if usr.DiscordID == "N/D" {
					continue
				}
				err := s.GuildMemberRoleAdd(ch.cfg.guildID, usr.DiscordID, ch.discord.tourneyRole.ID)
				if err != nil {
					log.Println(err.Error())
				}
			}
		} else {
			for _, usr := range ch.discord.contacts {
				if usr.DiscordID == "N/D" {
					continue
				}
				err := s.GuildMemberRoleRemove(ch.cfg.guildID, usr.DiscordID, ch.discord.tourneyRole.ID)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		embed = append(embed, ch.msgEmbed("Roles", []*discordgo.MessageEmbedField{
			{Name: "Success!"},
		}))
	} else {
		embed = append(embed, ch.msgEmbed("Roles", []*discordgo.MessageEmbedField{
			{Name: "Error: Can't work with roles by commands", Value: "CSV file with data isn't loaded. Load file and restart bot."},
		}))
	}
	return embed
}

// method for start sending messages for players tournament
func (ch *commandHandler) processSending(s *discordgo.Session, i *discordgo.InteractionCreate, local responseLocale) error {
	// Check values ID server (guildID) and URL to tournament (slug)
	if ch.cfg.guildID != "" && ch.slug != "" {
		if err := response(s, i, local.responseMsg.Starting); err != nil {
			return err
		}
		go ch.Process(s)
	}
	return fmt.Errorf("guildID = %v | slug = %v", ch.cfg.guildID, ch.slug)
}

// parse URL string for get slug value
func (ch *commandHandler) parseURL(i *discordgo.InteractionCreate, local responseLocale) ([]*discordgo.MessageEmbed, error) {
	embed := []*discordgo.MessageEmbed{}

	// parse string with URL
	u, err := url.Parse(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		embed = append(embed, ch.msgEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: "**Slug**", Value: local.errorMsg.Input},
		}))
		return embed, err
	}
	// separate URL to parts
	arg := strings.Split(u.Path, "/")

	// check URL on key words
	if len(arg) != 0 && arg[1] == "tournament" && arg[3] == "event" {
		ch.slug = arg[1] + "/" + arg[2] + "/" + arg[3] + "/" + arg[4]
		embed = append(embed, ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "**Slug**", Value: ch.slug},
		}))
		return embed, nil
	}

	embed = append(embed, ch.msgEmbed("Error", []*discordgo.MessageEmbedField{
		{Name: "**Slug**", Value: local.errorMsg.Input},
	}))
	return embed, fmt.Errorf(local.errorMsg.Input)
}

func (ch *commandHandler) getRuleMatchesData(i *discordgo.InteractionCreate, local responseLocale) []*discordgo.MessageEmbed {
	args := i.ApplicationCommandData().Options
	ch.cfg.rulesMatches.StandardFormat = int(args[0].IntValue())
	ch.cfg.rulesMatches.FinalsFormat = int(args[1].IntValue())
	ch.cfg.rulesMatches.Stage = args[2].StringValue()
	ch.cfg.rulesMatches.Rounds = int(args[3].IntValue())
	ch.cfg.rulesMatches.Duration = int(args[4].IntValue())
	ch.cfg.rulesMatches.Crossplatform = args[5].BoolValue()

	// Saving values in template msgRuleMatches
	return []*discordgo.MessageEmbed{ch.msgRuleMatches(i.Locale.String())}
}

func (ch *commandHandler) getStreamLobbyData(i *discordgo.InteractionCreate, local responseLocale) ([]*discordgo.MessageEmbed, error) {
	args := i.ApplicationCommandData().Options
	embed := []*discordgo.MessageEmbed{}

	code := regexp.MustCompile(`[0-9]+`).FindAllString(args[4].StringValue(), -1)[0]
	if len(code) != 4 {
		embed = append(embed, ch.msgEmbed("Error", []*discordgo.MessageEmbedField{
			{Name: local.vdMsg.MessageStreamHeader, Value: local.errorMsg.Input},
		}))
		return embed, fmt.Errorf("no 4 numbers in field")
	}

	ch.cfg.streamLobby.Area = args[0].StringValue()
	ch.cfg.streamLobby.Language = args[1].StringValue()
	ch.cfg.streamLobby.Conn = args[2].StringValue()
	ch.cfg.streamLobby.Crossplatform = args[3].BoolValue()
	ch.cfg.streamLobby.Passcode = code

	// Saving values in template msgStreamLobby
	embed = append(embed, ch.msgStreamLobby(i.Locale.String()))
	return embed, nil
}

func (ch *commandHandler) getLogoTournamnentURL(i *discordgo.InteractionCreate, local responseLocale) []*discordgo.MessageEmbed {
	arg := i.ApplicationCommandData().Options[0].StringValue()
	ch.cfg.tournament.Game.Name = arg

	return []*discordgo.MessageEmbed{ch.msgEmbed(local.vdMsg.LogoTournament, []*discordgo.MessageEmbedField{
		{Name: "**Url**", Value: fmt.Sprintf("%v", ch.cfg.tournament.Game.Name)},
	})}
}

func (ch *commandHandler) readCommandEmbedJSON(s *discordgo.Session, i *discordgo.InteractionCreate, local responseLocale) ([]*discordgo.MessageEmbed, error) {
	errRespond := func(embed []*discordgo.MessageEmbed) []*discordgo.MessageEmbed {
		embed = append(embed, ch.msgEmbed(local.vdMsg.Title, []*discordgo.MessageEmbedField{
			{Name: "", Value: local.errorMsg.NoData},
		}))

		return embed
	}
	embed := []*discordgo.MessageEmbed{}

	cts, err := os.ReadFile("contactsEmbed.json")
	if err != nil {
		return errRespond(embed), err
	}
	if err := json.Unmarshal(cts, &ch.discord.embedContacts); err != nil {
		return errRespond(embed), err
	}

	arg := i.ApplicationCommandData().Options[0].StringValue()
	switch strings.ToLower(arg) {
	case "any", "все":
		if err := response(s, i, local.responseMsg.InProcess); err != nil {
			return nil, err
		}
		for _, embedContact := range ch.discord.embedContacts {
			if _, err := s.ChannelMessageSendEmbed(i.ChannelID, embedContact); err != nil {
				return errRespond(embed), err
			}
		}
		return nil, nil
	default:
		for _, embedContact := range ch.discord.embedContacts {
			for _, field := range embedContact.Fields {
				// If argument from command == name from field which from json file
				if strings.EqualFold(arg, field.Name) {
					var fields []*discordgo.MessageEmbedField
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:  field.Name,
						Value: field.Value,
					})
					embed = append(embed, ch.msgEmbed(local.vdMsg.Title, fields))
					return embed, nil
				}
			}
		}
	}
	return errRespond(embed), fmt.Errorf("not finded player in json file: %v", arg)
}
