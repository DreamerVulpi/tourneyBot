package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigStartGG struct {
	Token string `toml:"token"`
}

type ConfigDiscordBot struct {
	Token   string `toml:"token"`
	GuildID string `toml:"guildID"`
	AppID   string `toml:"appID"`
}

type ConfigTempalate struct {
}

type Config struct {
	Startgg ConfigStartGG    `toml:"startgg"`
	Discord ConfigDiscordBot `toml:"discordbot"`
}

type Templates struct {
	InviteMessage string `toml:"invite"`
	StreamMessage string `toml:"stream"`
}

func LoadConfig(file string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	switch {
	case len(cfg.Startgg.Token) == 0:
		return Config{}, errors.New("startGG token is empty")
	case len(cfg.Discord.Token) == 0:
		return Config{}, errors.New("discord token is empty")
	case len(cfg.Discord.AppID) == 0:
		return Config{}, errors.New("discord appID is empty")
	case len(cfg.Discord.GuildID) == 0:
		return Config{}, errors.New("discord guildID is empty")
	default:
		return cfg, nil
	}
}

func LoadTemplates(file string) (Templates, error) {
	var tmpt Templates

	err := cleanenv.ReadConfig(file, &tmpt)
	if err != nil {
		return Templates{}, err
	}

	switch {
	case len(tmpt.InviteMessage) == 0:
		return Templates{}, errors.New("inviteMessage is empty")
	case len(tmpt.StreamMessage) == 0:
		return Templates{}, errors.New("streamMessage is empty")
	default:
		return tmpt, nil
	}

}
