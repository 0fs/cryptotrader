package main

import (
	"github.com/spf13/viper"
)

var config = viper.New()

func initConfig() {

	config.SetConfigFile("config/app.yml")

	err := config.ReadInConfig()

	if err != nil {
		logger.Fatal().Err(err).Msg("Fatal error config file")
	}
}
