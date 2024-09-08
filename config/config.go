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

type Config struct {
	Startgg ConfigStartGG    `toml:"startgg"`
	Discord ConfigDiscordBot `toml:"discordbot"`
}

type Template struct {
	InviteMessage string `toml:"invite"`
	StreamMessage string `toml:"stream"`
}

type ConfigTemplate struct {
	Template Template `toml:"template"`
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

func LoadTemplates(file string) (Template, error) {
	var t ConfigTemplate

	err := cleanenv.ReadConfig(file, &t)
	if err != nil {
		return Template{}, err
	}

	switch {
	case len(t.Template.InviteMessage) == 0:
		return Template{}, errors.New("inviteMessage is empty")
	case len(t.Template.StreamMessage) == 0:
		return Template{}, errors.New("streamMessage is empty")
	default:
		return t.Template, nil
	}

}
