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

func Load() {
	configFilename := "config.yaml"

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.upload")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file %s was not found.\n", configFilename)
		} else {
			panic(fmt.Errorf("Fatal error reading config file %s: error: %w\n",
				configFilename, err))
		}
	}
}
