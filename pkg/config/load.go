package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type OutputLocationType int

const (
	File OutputLocationType = iota
	Database
)

func (o OutputLocationType) String() string {
	return [...]string{"file", "database"}[o]
}

type OutputConfig struct {
	Type     string
	Location string
}

type Config struct {
	OriginalFile         string
	DataFile             string
	EnglishStopwordsFile string
	Output               OutputConfig
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
