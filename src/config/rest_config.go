package config

type RestConfig struct {
	Address string `env:"REST_ADDRESS,required"`
}
