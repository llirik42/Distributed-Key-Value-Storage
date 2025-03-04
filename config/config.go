package config

import (
	"errors"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type RaftConfig struct {
	BroadcastTimeMs      uint `env:"BROADCAST_TIME_MS,required"`
	MinElectionTimeoutMs uint `env:"MIN_ELECTION_TIMEOUT_MS,required"`
	MaxElectionTimeoutMs uint `env:"MAX_ELECTION_TIMEOUT_MS,required"`
	SelfNode             struct {
		Address string `env:"SELF_ADDRESS,required"`
		Id      string `env:"SELF_ID,required"`
	}
	OtherNodes []string `env:"OTHER_NODES,required"`
}

type Config struct {
	RaftConfig RaftConfig
}

func NewConfiguration() (*Config, error) {
	loadingErr := godotenv.Load()
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
