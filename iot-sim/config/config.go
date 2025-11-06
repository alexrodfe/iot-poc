// Package config provides configuration loading from ENV variables for the IoT Manager service.
package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Nats NatsConfig
}

type NatsConfig struct {
	URL string `mapstructure:"nats_url"`
}

func New() (*Config, error) {
	_ = godotenv.Load("../env") // mockup a real env

	v := viper.New()
	v.SetConfigFile("../.env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	v.AutomaticEnv()

	// NATS config
	_ = v.BindEnv("nats_url", "NATS_URL")

	natsConfig := NatsConfig{}
	if err := v.Unmarshal(&natsConfig); err != nil {
		return nil, err
	}

	config := Config{
		Nats: natsConfig,
	}

	return &config, nil
}
