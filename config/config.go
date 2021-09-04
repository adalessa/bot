package config

import "github.com/caarlos0/env/v6"

type Config struct {
	DiscordToken string `env:"TOKEN"`
	OpApiHost    string `env:"OP_API_HOST"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
