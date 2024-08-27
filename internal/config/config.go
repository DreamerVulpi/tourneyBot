package config

import "github.com/ilyakaznacheev/cleanenv"

type ConfigStartGG struct {
	Token string `toml:"token"`
}

type ConfigDiscordBot struct {
	Token string `toml:"token"`
}

type Config struct {
	Startgg ConfigStartGG    `toml:"startgg"`
	Discord ConfigDiscordBot `toml:"discordbot"`
}

func LoadConfig(file string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(file, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
