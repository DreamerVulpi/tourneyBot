package bot

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/config"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type commandHandler struct {
	slug         string
	guildID      string
	stop         chan struct{}
	m            *discordgo.MessageCreate
	client       *startgg.Client
	dataLobby    config.ConfigTournament
	RulesMatches config.RulesMatches
	StreamLobby  config.StreamLobby
	Bot          config.Bot
}

func response(s *discordgo.Session, i *discordgo.InteractionCreate, text string) error {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: text,
			},
		},
	)
	if err != nil {
		return errors.New("can't respond on message")
	}
	return nil
}

func (cmd *commandHandler) messageEmbed(title string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: title,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: cmd.Bot.Img,
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cmd.Bot.LogoTournament,
		},
	}
}

func (cmd *commandHandler) responseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed []*discordgo.MessageEmbed) error {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embed,
			},
		},
	)
	return err
}

func (cmd *commandHandler) check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, cmd.messageEmbed("Check data", []*discordgo.MessageEmbedField{
		{Name: "**Slug**", Value: fmt.Sprintln("A slug is made of two parts, the tournament name and the event name. The format is this:\n*tournament/<tournament_name>/event/<event_name>*"), Inline: true},
		{Value: fmt.Sprintf("```%v```", cmd.slug)},
	}))

	embed = append(embed, cmd.messageEmbed("Rules matches", []*discordgo.MessageEmbedField{
		{Name: "**Format**", Value: fmt.Sprintf("FT%v", cmd.RulesMatches.Format), Inline: true},
		{Name: "**Stage**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Stage), Inline: true},
		{Name: "**Rounds in 1 match**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Rounds)},
		{Name: "**Seconds in 1 round**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Duration), Inline: true},
		{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Crossplatform), Inline: true},
	}))

	embed = append(embed, cmd.messageEmbed("Stream lobby data", []*discordgo.MessageEmbedField{
		{Name: "**Area**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Area), Inline: true},
		{Name: "**Language**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Language), Inline: true},
		{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Crossplatform)},
		{Name: "**Passcode**", Value: fmt.Sprintf("```%v```", cmd.StreamLobby.Passcode), Inline: true},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("check: can't respond on message"))
	}
}
func (cmd *commandHandler) start_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := response(s, i, "Start sending..."); err != nil {
		log.Println(err.Error())
	}

	go cmd.SendingMessages(s)
}
func (cmd *commandHandler) stop_sending(s *discordgo.Session, i *discordgo.InteractionCreate) {
	go func() {
		response(s, i, "stopping...")
	}()

	// Send signal to stop process
	cmd.stop <- struct{}{}

	s.ChannelMessageSend(i.ChannelID, "Stopped!")
}

// TODO: PARSE LINK ON EVENT USING PACKAGE "URL"
func (cmd *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.slug = i.ApplicationCommandData().Options[0].StringValue()

	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, cmd.messageEmbed("Check data", []*discordgo.MessageEmbedField{
		{Name: "**Slug**", Value: cmd.slug},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("setEvent: can't respond on message"))
	}
}

func (cmd *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options

	format := int(args[0].IntValue())
	if format <= 10 {
		cmd.RulesMatches.Format = format
	}

	rounds := int(args[2].IntValue())
	if rounds <= 5 {
		cmd.RulesMatches.Rounds = rounds
	}

	duration := int(args[3].IntValue())
	if duration <= 99 {
		cmd.RulesMatches.Duration = duration
	}

	cmd.RulesMatches.Stage = args[1].StringValue()
	cmd.RulesMatches.Crossplatform = args[4].BoolValue()

	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, cmd.messageEmbed("Check data", []*discordgo.MessageEmbedField{
		{Name: "**Rules matches**", Value: ""},
		{Name: "**Format**", Value: fmt.Sprintf("FT%v", cmd.RulesMatches.Format) + fmt.Sprintf(" (First to %v win in set)", cmd.RulesMatches.Format), Inline: true},
		{Name: "**Stage**", Value: cmd.RulesMatches.Stage, Inline: true},
		{Name: "**Rounds in 1 match**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Rounds)},
		{Name: "**Time in 1 round**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Duration) + " seconds"},
		{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Crossplatform)},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("editRuleMatches: can't respond on message"))
	}
}

func (cmd *commandHandler) editStreamLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options

	area := args[0].StringValue()
	if len(area) != 0 {
		cmd.StreamLobby.Area = area
	}

	lang := args[1].StringValue()
	if len(lang) != 0 {
		cmd.StreamLobby.Language = lang
	}

	conn := args[2].StringValue()
	if len(conn) != 0 {
		cmd.StreamLobby.Conn = args[2].StringValue()
	}

	cmd.StreamLobby.Crossplatform = args[3].BoolValue()

	passcode := args[4].StringValue()
	if len(passcode) <= 4 {
		cmd.StreamLobby.Passcode = passcode
	}

	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, cmd.messageEmbed("Stream lobby", []*discordgo.MessageEmbedField{
		{Name: "**Area**", Value: fmt.Sprintf("FT%v", cmd.StreamLobby.Area)},
		{Name: "**Language**", Value: cmd.RulesMatches.Stage},
		{Name: "**Connection quality preference**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Conn)},
		{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Crossplatform)},
		{Name: "**Passcode**", Value: fmt.Sprintf("%v", cmd.StreamLobby.Passcode)},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("editStreamLobby: can't respond on message"))
	}
}

func (cmd *commandHandler) editLogoTournament(s *discordgo.Session, i *discordgo.InteractionCreate) {
	arg := i.ApplicationCommandData().Options[0].StringValue()
	cmd.Bot.LogoTournament = arg

	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, cmd.messageEmbed("Logo tournament", []*discordgo.MessageEmbedField{
		{Name: "**Url**", Value: fmt.Sprintf("%v", cmd.Bot.LogoTournament)},
	}))

	if err := cmd.responseEmbed(s, i, embed); err != nil {
		log.Println(errors.New("editLogoTournament: can't respond on message"))
	}
}

// TODO: Add new command: Help
