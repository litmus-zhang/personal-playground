package config

type Config struct {
	DbDriver          string `mapstructure:"DB_DRIVER"`
	DbSource          string `mapstructure:"DB_SOURCE"`
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}
