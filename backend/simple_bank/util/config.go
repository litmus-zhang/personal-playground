package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER" envconfig:"DB_DRIVER" default:"postgres"`
	DBURL         string `mapstructure:"DB_URL" envconfig:"DB_URL" default:"postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" envconfig:"SERVER_ADDRESS" default:":8080"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
