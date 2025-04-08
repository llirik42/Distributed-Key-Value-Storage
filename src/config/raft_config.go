package config

type RaftConfig struct {
	BroadcastTimeMs      int `env:"BROADCAST_TIME_MS,required"`
	MinElectionTimeoutMs int `env:"MIN_ELECTION_TIMEOUT_MS,required"`
	MaxElectionTimeoutMs int `env:"MAX_ELECTION_TIMEOUT_MS,required"`
	SelfNode             struct {
		Address string `env:"SELF_ADDRESS,required"`
		Id      string `env:"SELF_ID,required"`
	}
	OtherNodes          []string `env:"OTHER_NODES,required"`
	ExecutedCommandsKey string   `env:"EXECUTED_COMMANDS_KEY,required"`
}
