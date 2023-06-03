package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(cfg string, data interface{}) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bingoctl")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		Cfg = NewDefaultConfig()

		return
	}

	if err := viper.Unmarshal(data); err != nil {
		fmt.Println(err)
	}
}
