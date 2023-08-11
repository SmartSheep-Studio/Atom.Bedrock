package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewConfigProvider(log zerolog.Logger) *viper.Viper {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading config")
	}

	return viper.GetViper()
}
