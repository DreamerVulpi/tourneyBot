package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (ch *commandHandler) workRoles(s *discordgo.Session, arg string) []*discordgo.MessageEmbed {
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
