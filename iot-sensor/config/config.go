// Package config provides configuration loading from ENV variables for the IoT Manager service.
package config

import (
	"github.com/alexrodfe/iot-poc/commons"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	SensorID string
	Nats     NatsConfig
}

type NatsConfig struct {
	URL                string
	SensorStream       string
	MeasurementsStream string
}

func New(sensorID string) (*Config, error) {
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

	natsConfig := NewNatsConfig(sensorID, v.GetString("nats_url"))

	config := Config{
		SensorID: sensorID,
		Nats:     natsConfig,
	}

	return &config, nil
}

func NewNatsConfig(sensorID, url string) NatsConfig {
	return NatsConfig{
		URL:                url,
		SensorStream:       commons.GetSensorStream(sensorID),
		MeasurementsStream: commons.GetMeasurementsStream(sensorID),
	}
}
