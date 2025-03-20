package config

type RESTConfig struct {
	Address string `env:"REST_ADDRESS,required"`
}
