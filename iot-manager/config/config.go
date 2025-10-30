// Package config provides configuration loading from ENV variables for the IoT Manager service.
package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Nats NatsConfig `mapstructure:"nats"`
}

type NatsConfig struct {
	URL         string `mapstructure:"url"`
	AdminStream string `mapstructure:"admin_stream"`
}

func New() (*Config, error) {
	_ = godotenv.Load(".env") // mockup a real env

	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	v.SetEnvPrefix("MANAGER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// NATS config
	_ = v.BindEnv("nats.url", "NATS_URL")
	_ = v.BindEnv("nats.admin_stream", "NATS_ADMIN_STREAM")

	cfg := Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
