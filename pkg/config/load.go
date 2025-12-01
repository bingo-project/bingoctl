package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(cfg string, data interface{}) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		// Get User home dir
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Add `$HOME` & `.`
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(".bingo")
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
