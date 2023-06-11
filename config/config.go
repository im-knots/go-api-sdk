package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func NewConfig(configPath string) *Config {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return &Config{viper: v}
}

func (c *Config) Unmarshal(output interface{}) error {
	return c.viper.Unmarshal(output)
}
