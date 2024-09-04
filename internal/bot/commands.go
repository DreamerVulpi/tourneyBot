package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "check",
			Description:              "Check variables startgg, discord and bot.",
			DMPermission:             &dmPermission,
			DefaultMemberPermissions: &defaultMemberPermissions,
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "проверка",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Проверить переменные startgg, discord, и бота.",
			},
		},
		{
			Name:        "set-event",
			Description: "Set event in configuration bot for getting all phaseGroups",
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "slug",
				Description: "Format: tournament/<tournament_name>/event/<event_name>",
				Required:    true,
			}},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "установить-ивент",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Установить идентификатор ивента в конфигурацию бота для получения всех групп",
			},
		},
		{
			Name:        "set-server-id",
			Description: "Set guild ID in configuration bot for getting members server",
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "slug",
				Description: "Guild ID = Server ID",
				Required:    true,
			}},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "установить-идентификатор-сервера",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Установить guild ID в конфигурацию бота для получения доступа к списку пользователей",
			},
		},
		{
			Name:        "start-sending",
			Description: "Start sending invite-messages to members tournament",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "начать-рассылку",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Начать рассылку инвайт-сообщений участникам турнира",
			},
		},
		{
			Name:        "stop-sending",
			Description: "Stop sending invite-messages to tournament sets",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "остановить-рассылку",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Остановить рассылку инвайт-сообщений участникам турнира",
			},
		},

		{
			Name:        "edit-invite-message",
			Description: "Edit template invite message",
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "template",
				Description: "edited template",
				Required:    true,
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Russian: "Отредактированный шаблон",
				},
			}},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "редактировать-инвайт-сообщения",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Редактировать шаблон инвайт-сообщения",
			},
		},
	}
)
