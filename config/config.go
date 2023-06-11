package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func NewConfig(configPath string) *Config {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore since we can use env vars
			log.Println("No config file found. Relying on environment variables.")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Fatal error config file: %s", err)
		}
	}

	return &Config{viper: v}
}

func (c *Config) Unmarshal(output interface{}) error {
	return c.viper.Unmarshal(output)
}
