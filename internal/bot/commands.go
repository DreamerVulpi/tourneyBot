package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func SetCommands(appID, guildID string, session *discordgo.Session) (*discordgo.Session, error) {
	if !app() {
		return &discordgo.Session{}, errors.New("app ID isn't correct or empty")
	}
	if !server() {
		return &discordgo.Session{}, errors.New("server ID (Guild ID) isn't correct or empty")
	}

	_, err := session.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{
		{
			Name:        "check",
			Description: "Check status startgg, discord and bot",
		},
		{
			Name:        "set-event",
			Description: "Set event in configuration bot for getting all phaseGroups",
		},
		{
			Name:        "set-server-id",
			Description: "Set guard ID for sending messages to members server",
		},
		{
			Name:        "start-sending",
			Description: "Start sending invite-messages to tournament sets",
		},
		{
			Name:        "stop-sending",
			Description: "Stop sending invite-messages to tournament sets",
		},

		{
			Name:        "edit-invite-message",
			Description: "Edit template invite message",
		},
	})
	return session, err
}
