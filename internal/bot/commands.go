package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	// TODO: Change access commands to only administrator server or specic role
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
			Name:        "edit-rules",
			Description: "Edit rule matches",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "format",
					Description: "First to ? wins",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "До ? побед",
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "map",
					Description: "Name map | Random (Example: yakushima | any and etc)",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Название карты | Любая (Ответ только по английски. К примеру: yakushima | any и т.д)",
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "rounds",
					Description: "Rounds in 1 match",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Раундов в 1 матче",
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Seconds in 1 round",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Секунд в 1 раунде",
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "crossplatformplay",
					Description: "Cross-platform game support (true | false)",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Поддержка кроссплатформенной игры (true | false)",
					},
				},
			},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "редактировать-правила-матчей",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Редактировать правила матчей в сете",
			},
		},
	}
)
