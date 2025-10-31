// Package config provides configuration loading from ENV variables for the IoT Manager service.
package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	SensorID uint64     `mapstructure:"sensor_id"`
	Nats     NatsConfig `mapstructure:"nats"`
}

type NatsConfig struct {
	URL                string `mapstructure:"url"`
	SensorStream       string `mapstructure:"sensor_stream"`
	MeasurementsStream string `mapstructure:"measurements_stream"`
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

	v.AutomaticEnv()

	// Sensor config
	_ = v.BindEnv("sensor_id", "SENSOR_ID")

	// NATS config
	_ = v.BindEnv("nats.url", "NATS_URL")
	_ = v.BindEnv("nats.sensor_stream", "NATS_SENSOR_STREAM")
	_ = v.BindEnv("nats.measurements_stream", "NATS_MEASUREMENTS_STREAM")

	cfg := Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
