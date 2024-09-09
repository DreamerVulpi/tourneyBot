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

type Local struct {
	Area      string `toml:"area"`
	Language  string `toml:"language"`
	Conn      string `toml:"connection"`
	Victory   int    `toml:"victory"`
	Rounds    int    `toml:"rounds"`
	Duration  int    `toml:"duration"`
	CrossPlay string `toml:"crossplay"`
}

type Stream struct {
	Area      string `toml:"area"`
	Language  string `toml:"language"`
	CrossPlay string `toml:"crossplay"`
	Passcode  string `toml:"passcode"`
}

type ConfigLobby struct {
	Local  Local  `toml:"local"`
	Stream Stream `toml:"stream"`
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

func LoadLobby(file string) (ConfigLobby, error) {
	var l ConfigLobby

	err := cleanenv.ReadConfig(file, &l)
	if err != nil {
		return ConfigLobby{}, err
	}

	switch {
	case len(l.Local.Area) == 0:
		return ConfigLobby{}, errors.New("local: field area is empty")
	case len(l.Local.Language) == 0:
		return ConfigLobby{}, errors.New("local: field language is empty")
	case len(l.Local.Conn) == 0:
		return ConfigLobby{}, errors.New("local: field connection is empty")
	case l.Local.Victory == 0:
		return ConfigLobby{}, errors.New("local: field victory is null")
	case l.Local.Rounds == 0:
		return ConfigLobby{}, errors.New("local: field rounds is empty")
	case l.Local.Duration == 0:
		return ConfigLobby{}, errors.New("local: field duration is empty")
	case len(l.Local.CrossPlay) == 0:
		return ConfigLobby{}, errors.New("local: field crossplay is empty")
	case len(l.Stream.Area) == 0:
		return ConfigLobby{}, errors.New("stream: field area is empty")
	case len(l.Stream.Language) == 0:
		return ConfigLobby{}, errors.New("stream: field language is empty")
	case len(l.Stream.CrossPlay) == 0:
		return ConfigLobby{}, errors.New("stream: field crossplay is empty")
	case len(l.Stream.Passcode) == 0:
		return ConfigLobby{}, errors.New("stream: field passcode is empty")
	default:
		return l, nil
	}
}
