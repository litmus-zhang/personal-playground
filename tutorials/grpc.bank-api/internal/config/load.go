package config

import (
	"log"

	"github.com/spf13/viper"
)

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("secret")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	log.Printf("config: %+v", config)
	return
}
