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

	viper.SetConfigName("upload-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/mnt/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
      fmt.Printf("upload: Config file was not found.\n")
		} else {
			panic(fmt.Errorf("Fatal error reading config file: error: %w\n", err))
		}
	}
}
