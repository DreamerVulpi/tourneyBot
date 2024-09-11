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

type RulesMatches struct {
	Format        int    `toml:"format"`
	Stage         string `toml:"stage"`
	Rounds        int    `toml:"rounds"`
	Duration      int    `toml:"duration"`
	Crossplatform bool   `toml:"crossplatform"`
}

type StreamLobby struct {
	Area          string `toml:"area"`
	Language      string `toml:"language"`
	Conn          string `toml:"connection"`
	Crossplatform bool   `toml:"crossplatform"`
	Passcode      string `toml:"passcode"`
}

type Bot struct {
	Img            string `toml:"img"`
	LogoTournament string `toml:"logoTournament"`
}

type ConfigTournament struct {
	Rules  RulesMatches `toml:"rules"`
	Stream StreamLobby  `toml:"stream"`
	Bot    Bot          `toml:"bot"`
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

func LoadTournament(file string) (ConfigTournament, error) {
	var l ConfigTournament

	err := cleanenv.ReadConfig(file, &l)
	if err != nil {
		return ConfigTournament{}, err
	}

	switch {
	case l.Rules.Format == 0:
		return ConfigTournament{}, errors.New("local: field format is null")
	case l.Rules.Rounds == 0:
		return ConfigTournament{}, errors.New("local: field rounds is empty")
	case len(l.Rules.Stage) == 0:
		return ConfigTournament{}, errors.New("local: field stage is empty")
	case l.Rules.Duration == 0:
		return ConfigTournament{}, errors.New("local: field duration is empty")
	case len(l.Stream.Area) == 0:
		return ConfigTournament{}, errors.New("stream: field area is empty")
	case len(l.Stream.Language) == 0:
		return ConfigTournament{}, errors.New("stream: field language is empty")
	case len(l.Stream.Conn) == 0:
		return ConfigTournament{}, errors.New("stream: field connection is empty")
	case len(l.Stream.Passcode) == 0:
		return ConfigTournament{}, errors.New("stream: field passcode is empty")
	case len(l.Bot.Img) == 0:
		return ConfigTournament{}, errors.New("bot: field img is empty")
	case len(l.Bot.LogoTournament) == 0:
		return ConfigTournament{}, errors.New("bot: field logoTournament is empty")
	default:
		return l, nil
	}
}
