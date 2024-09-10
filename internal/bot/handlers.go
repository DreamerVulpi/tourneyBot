package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
	dataLobby    config.ConfigLobby
	RulesMatches config.RulesMatches
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

func responseSetted(s *discordgo.Session, i *discordgo.InteractionCreate, msgformat string, margs []interface{}) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
	if err != nil {
		return errors.New("can't respond on message")
	}
	return nil
}

func (cmd *commandHandler) check(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := []*discordgo.MessageEmbed{}
	embed = append(embed, &discordgo.MessageEmbed{
		Title: "Check data",
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: "https://i.imgur.com/AfFp7pu.png",
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "**Slug**", Value: fmt.Sprintln("A slug is made of two parts, the tournament name and the event name. The format is this:```tournament/<tournament_name>/event/<event_name>```"), Inline: true},
			{Value: cmd.slug},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/AfFp7pu.png",
		},
	})

	embed = append(embed, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: "https://i.imgur.com/AfFp7pu.png",
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "**Rules matches**"},
			{Name: "**Format**", Value: fmt.Sprintf("FT%v", cmd.RulesMatches.Format)},
			{Name: "**Map**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Stage)},
			{Name: "**Rounds in 1 match**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Rounds)},
			{Name: "**Seconds in 1 round**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Duration)},
			{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.RulesMatches.Crossplatform)},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/AfFp7pu.png",
		},
	})

	embed = append(embed, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: "https://i.imgur.com/AfFp7pu.png",
			URL:     "https://github.com/DreamerVulpi/tourneybot",
			Name:    "TourneyBot",
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "**Stream lobby data**"},
			{Name: "**Area**", Value: fmt.Sprintf("%v", cmd.dataLobby.Stream.Area)},
			{Name: "**Language**", Value: fmt.Sprintf("%v", cmd.dataLobby.Stream.Language)},
			{Name: "**Crossplatform**", Value: fmt.Sprintf("%v", cmd.dataLobby.Stream.Crossplatform)},
			{Name: "**Passcode**", Value: fmt.Sprintf("```%v```", cmd.dataLobby.Stream.Passcode)},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/AfFp7pu.png",
		},
	})

	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embed,
			},
		},
	)
	if err != nil {
		log.Println(errors.New("can't respond on message"))
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

func (cmd *commandHandler) setEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()
	cmd.slug = input

	margs := make([]interface{}, 0, len(input))
	msgformat := ""

	margs = append(margs, input)
	msgformat += "> SLUG: %s\n"

	if err := responseSetted(s, i, msgformat, margs); err != nil {
		log.Println(err.Error())
	}
}

// TODO: Refactor code
func (cmd *commandHandler) editRuleMatches(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := i.ApplicationCommandData().Options

	margs := make([]interface{}, 0, len(args))
	msgformat := ""

	format := int(args[0].IntValue())
	if format <= 10 {
		cmd.RulesMatches.Format = format
		margs = append(margs, format)
	}

	listStages := map[string]string{
		"arena":                         "Arena",
		"arena (underground)":           "Arena (Underground)",
		"urban square":                  "Urban Square",
		"urban square (evening)":        "Urban Square (Evening)",
		"yakushima":                     "Yakushima",
		"coliseum of fate":              "Coliseum of Fate",
		"rebel hangar":                  "Rebel Hanger",
		"fallen destiny":                "Fallen Destiny",
		"descent into subconsciousness": "Descent into Subconsciousness",
		"sanctum":                       "Sanctum",
		"into the stratpsphere":         "Into the Stratosphere",
		"ortiz farm":                    "Ortiz Farm",
		"celebration on the seine":      "Celebration on the Seine",
		"secluded training ground":      "Secluded Training Ground",
		"elegant palace":                "Elegant Palace",
		"midnight siege":                "Midnight Siege",
		"seaside resort":                "Seaside Resort",
		"any":                           "Any",
	}

	stage := args[1].StringValue()
	if len(listStages[strings.ToLower(stage)]) != 0 {
		cmd.RulesMatches.Stage = listStages[strings.ToLower(stage)]
		margs = append(margs, stage)
	}

	rounds := int(args[2].IntValue())
	if rounds <= 5 {
		cmd.RulesMatches.Rounds = rounds
		margs = append(margs, rounds)
	}

	duration := int(args[3].IntValue())
	if duration <= 99 {
		cmd.RulesMatches.Duration = duration
		margs = append(margs, duration)
	}

	cmd.RulesMatches.Crossplatform = args[4].BoolValue()
	margs = append(margs, cmd.RulesMatches.Crossplatform)

	msgformat += "> Saved data: %s\n"

	if err := responseSetted(s, i, msgformat, margs); err != nil {
		log.Println(err.Error())
	}
}

// TODO: Add new command: urlLogo
// TODO: Add new command: editStreamMessage (with more args)
// TODO: Add new command: Help
// TODO: Add new command: About
