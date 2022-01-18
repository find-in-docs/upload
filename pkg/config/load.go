package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OriginalFile         string
	DataDir              string
	EnglishStopwordsFile string
	OutputDir            string
}

func LoadConfig() *Config {
	configFilename := "config.yaml"

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.search")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file %s was not found.\n", configFilename)
		} else {
			panic(fmt.Errorf("Fatal error reading config file %s: error: %w\n",
				configFilename, err))
		}
	}

	var C Config
	if err := viper.Unmarshal(&C); err != nil {
		panic(fmt.Errorf("Fatal error unmarshaling config file %s.\n", configFilename))
	}

	return &C
}
