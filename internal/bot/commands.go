package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneybot/config"
)

func (cmd *commandHandler) commands() []*discordgo.ApplicationCommand {
	// TODO: Change access commands to only administrator server or specic role
	return []*discordgo.ApplicationCommand{
		{
			Name:        "check",
			Description: "Check variables startgg, discord and bot.",
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
				Name:        "link",
				Description: "Link on event must including path: tournament/<tournament_name>/event/<event_name>",
				Required:    true,
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Russian: "Cсылка на ивент должна включать в себя путь: tournament/<название_турнира>/event/<название_ивента>",
				},
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
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "First to 1 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 1 победы в сете",
							},
							Value: 1,
						},
						{
							Name: "First to 2 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 2 побед в сете",
							},
							Value: 2,
						},
						{
							Name: "First to 3 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 3 побед в сете",
							},
							Value: 3,
						},
						{
							Name: "First to 4 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 4 побед в сете",
							},
							Value: 4,
						},
						{
							Name: "First to 5 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 5 побед в сете",
							},
							Value: 5,
						},
						{
							Name: "First to 6 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 6 побед в сете",
							},
							Value: 6,
						},
						{
							Name: "First to 7 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 7 побед в сете",
							},
							Value: 7,
						},
						{
							Name: "First to 8 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 8 побед в сете",
							},
							Value: 8,
						},
						{
							Name: "First to 9 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 9 побед в сете",
							},
							Value: 9,
						},
						{
							Name: "First to 10 win in set",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "До 10 побед в сете",
							},
							Value: 10,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "stage",
					Description: "Name stage | Random",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Название арены | Любая",
					},
					Choices: append(choice(config.ListStages), &discordgo.ApplicationCommandOptionChoice{
						Name: "Any",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.Russian: "Любая",
						},
						Value: "any",
					}),
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "rounds",
					Description: "Rounds in 1 match",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Раундов в 1 матче",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "1",
							Value: 1,
						},
						{
							Name:  "2",
							Value: 2,
						},
						{
							Name:  "3",
							Value: 3,
						},
						{
							Name:  "4",
							Value: 4,
						},
						{
							Name:  "5",
							Value: 5,
						},
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
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "30",
							Value: 30,
						},
						{
							Name:  "40",
							Value: 40,
						},
						{
							Name:  "60",
							Value: 60,
						},
						{
							Name:  "80",
							Value: 80,
						},
						{
							Name:  "99",
							Value: 99,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "crossplatformplay",
					Description: "Cross-platform game support",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Поддержка кроссплатформенной игры",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Enable",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Включено",
							},
							Value: true,
						},
						{
							Name: "Disable",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Выключено",
							},
							Value: false,
						},
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
		{
			Name:        "edit-stream-lobby",
			Description: "Edit stream lobby configuration",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "area",
					Description: "Area",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Регион",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Any",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Любой",
							},
							Value: "any",
						},
						{
							Name: "Close to Me",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Ближе ко мне",
							},
							Value: "close",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "language",
					Description: "Language",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Язык",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Any",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Любой",
							},
							Value: "any",
						},
						{
							Name: "Same as Me",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Как у меня",
							},
							Value: "same",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "conn",
					Description: "Connection quality preference",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Качество соединения",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "No Restrictions",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Нет ограничений",
							},
							Value: "any",
						},
						{
							Name: "3 or better",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "3 и больше",
							},
							Value: "3",
						},
						{
							Name: "4 or better",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "4 и больше",
							},
							Value: "4",
						},
						{
							Name: "5 or better",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "5",
							},
							Value: "5",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "crossplatformplay",
					Description: "Cross-platform game support",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Поддержка кроссплатформенной игры",
					},
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Enable",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Включено",
							},
							Value: true,
						},
						{
							Name: "Disable",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian: "Выключено",
							},
							Value: false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "passcode",
					Description: "Passcode to access your stream lobby (min: 0000; max: 9999)",
					Required:    true,
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian: "Пароль к лобби для стрима (мин: 0000; макс:9999)",
					},
				},
			},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "редактировать-стрим-лобби",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Редактировать конфигурацию лобби для стрима",
			},
		},
		{
			Name:        "edit-logo-tournament",
			Description: "Edit link to logo tournament",
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "Link to logo tournament",
				Required:    true,
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Russian: "Cсылка на логотип турнира",
				},
			}},
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "редактировать-лого-турнира",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Редактировать ссылку на логотип турнира",
			},
		},
	}
}

func choice(list map[string]string) []*discordgo.ApplicationCommandOptionChoice {
	var result []*discordgo.ApplicationCommandOptionChoice
	for key, value := range list {
		result = append(result, &discordgo.ApplicationCommandOptionChoice{
			Name:  value,
			Value: key,
		})
	}
	return result
}
