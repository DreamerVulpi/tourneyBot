package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// FIXME: NOT WORKING!>>#@$@>#$

func handlerInputs(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	m.Content = strings.TrimSpace(m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "/set-event") {
		setEvent(s, m, args[1])
	}

	if strings.HasPrefix(m.Content, "/set-server-id") {
		setEvent(s, m, args[1])
	}

	if strings.HasPrefix(m.Content, "/edit-invite-message") {
		setEvent(s, m, args[1])
	}

}

func setEvent(s *discordgo.Session, m *discordgo.MessageCreate, data string) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	SetSlug(data)

	_, err = s.ChannelMessageSend(channel.ID, "Success saved!")
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

func editInviteMessage(s *discordgo.Session, m *discordgo.MessageCreate, data string) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	templateInviteMessage = data

	_, err = s.ChannelMessageSend(channel.ID, "Success saved!")
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

func setServerID(s *discordgo.Session, m *discordgo.MessageCreate, data string) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	SetServerID(data)

	_, err = s.ChannelMessageSend(channel.ID, "Success saved!")
	if err != nil {
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}
