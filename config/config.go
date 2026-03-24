package config

import (
	"errors"
	"log"

	"path/filepath"

	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDiscordBot struct {
	Token          string `toml:"token"`
	GuildID        string `toml:"guildID"`
	DebugChannelID string `toml:"debugChannelID"`
}

type ConfigGame struct {
	Name string `toml:"name"`
}

type ConfigRolesIdDiscord struct {
	Ru string `toml:"ru"`
}

type DebugMode struct {
	Mode bool `toml:"mode"`
}

type Database struct {
	Dsn string `toml:"dsn"`
}

type Config struct {
	Discord   ConfigDiscordBot     `toml:"discordbot"`
	Roles     ConfigRolesIdDiscord `toml:"roles"`
	DebugMode DebugMode            `toml:"debug"`
	Db        Database             `toml:"database"`
}

type RulesMatches struct {
	StandardFormat int    `toml:"standardFormat"`
	FinalsFormat   int    `toml:"finalsFormat"`
	Stage          string `toml:"stage"`
	Rounds         int    `toml:"rounds"`
	Duration       int    `toml:"duration"`
	Crossplatform  bool   `toml:"crossplatform"`
	Waiting        int    `toml:"waiting"`
}

type StreamLobby struct {
	Area          string `toml:"area"`
	Language      string `toml:"language"`
	Conn          string `toml:"connection"`
	Crossplatform bool   `toml:"crossplatform"`
	Passcode      string `toml:"passcode"`
}

type Logo struct {
	Img string `toml:"img"`
}

type Csv struct {
	NameFile string `toml:"nameFile"`
}

type NamePlatform struct {
	Platform string `toml:"platform"`
}

type ConfigTournament struct {
	Platform NamePlatform `toml:"tournamentPlatform"`
	Rules    RulesMatches `toml:"rules"`
	Stream   StreamLobby  `toml:"stream"`
	Logo     Logo         `toml:"logo"`
	Csv      Csv          `toml:"csv"`
	Game     ConfigGame   `toml:"game"`
}

func GetAbsPath(relativeToRoot string) string {
	// path file
	_, filename, _, _ := runtime.Caller(0)
	// ../
	root := filepath.Join(filepath.Dir(filename), "..")

	absPath := filepath.Join(root, relativeToRoot)
	return absPath
}

func LoadConfig(file string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	switch {
	case len(cfg.Discord.Token) == 0:
		return Config{}, errors.New("discord: token is empty")
	case len(cfg.Discord.GuildID) == 0:
		return Config{}, errors.New("discord: guildID is empty")
	case len(cfg.Discord.DebugChannelID) == 0:
		return Config{}, errors.New("discord: debugChannelID is empty")
	case len(cfg.Roles.Ru) == 0:
		log.Println(errors.New("roles: ru locale is empty").Error())
		return cfg, nil
	case len(cfg.Db.Dsn) == 0:
		return Config{}, errors.New("postgres: dsn string is empty")
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
	case len(l.Platform.Platform) == 0:
		return ConfigTournament{}, errors.New("local: field platform is null")
	case l.Rules.StandardFormat == 0:
		return ConfigTournament{}, errors.New("local: field standardFormat is null")
	case l.Rules.FinalsFormat == 0:
		return ConfigTournament{}, errors.New("local: field finalsFormat is null")
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
	case l.Rules.Waiting > 30 || l.Rules.Waiting <= 0:
		return ConfigTournament{}, errors.New("waiting time: isn't correct")
	case len(l.Game.Name) == 0:
		return ConfigTournament{}, errors.New("game: name is empty")
	case len(l.Csv.NameFile) == 0:
		log.Println(errors.New("csv: nameFile field is empty").Error())
		return l, nil
	case len(l.Logo.Img) == 0:
		log.Println(errors.New("tournament: logo link is empty").Error())
		return l, nil
	default:
		return l, nil
	}
}
