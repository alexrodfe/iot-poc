// Package config provides configuration loading from ENV variables for the IoT Manager service.
package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	SensorID string `mapstructure:"sensor_id"`
	Nats     NatsConfig
}

type NatsConfig struct {
	URL                 string `mapstructure:"nats_url"`
	SensorStream        string `mapstructure:"nats_sensor_stream"`
	SensorSubject       string `mapstructure:"nats_sensor_subject"`
	MeasurementsStream  string `mapstructure:"nats_measurements_stream"`
	MeasurementsSubject string `mapstructure:"nats_measurements_subject"`
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

	// Sensor config
	_ = v.BindEnv("sensor_id", "SENSOR_ID")

	// NATS config
	_ = v.BindEnv("nats_url", "NATS_URL")
	_ = v.BindEnv("nats_sensor_stream", "NATS_SENSOR_STREAM")
	_ = v.BindEnv("nats_sensor_subject", "NATS_SENSOR_SUBJECT")
	_ = v.BindEnv("nats_measurements_stream", "NATS_MEASUREMENTS_STREAM")
	_ = v.BindEnv("nats_measurements_subject", "NATS_MEASUREMENTS_SUBJECT")

	natsConfig := NatsConfig{}
	if err := v.Unmarshal(&natsConfig); err != nil {
		return nil, err
	}

	config := Config{
		SensorID: v.GetString("sensor_id"),
		Nats:     natsConfig,
	}

	return &config, nil
}
