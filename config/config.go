package config

import (
	"errors"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	RaftConfig RaftConfig
}

func NewConfiguration(filePath string) (*Config, error) {
	loadingErr := godotenv.Load(filePath)
	if loadingErr != nil {
		return nil, errors.Join(errors.New("error reading configuration"), loadingErr)
	}

	cfg := &Config{}
	parsingErr := env.Parse(cfg)

	if parsingErr != nil {
		return nil, errors.Join(errors.New("error parsing configuration"), parsingErr)
	}
	return cfg, nil
}
