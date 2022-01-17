package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OriginalFile string
	DataDir      string
	OutputDir    string
}

// Why is this function not considered declared when this package is imported into main.go?
func LoadConfig() *Config {
	config_filename := "config.yaml"

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.search")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file %s was not found.\n", config_filename)
		} else {
			panic(fmt.Errorf("Fatal error reading config file %s: error: %w\n",
				config_filename, err))
		}
	}

	var C Config
	if err := viper.Unmarshal(&C); err != nil {
		panic(fmt.Errorf("Fatal error unmarshaling config file %s.\n", config_filename))
	}

	return &C
}
