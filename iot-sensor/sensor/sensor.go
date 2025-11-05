// Package sensor declares the IoT sensor
package sensor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sensor/clients"
	"github.com/alexrodfe/iot-poc/iot-sensor/config"
)

type Sensor struct {
	ID string

	sensorConfig sensorConfig
	natsConfig   config.NatsConfig

	natsClient clients.Nats

	generatedMetrics map[uint64]commons.Measurement
}

func New(cfg *config.Config, natsClient clients.Nats) *Sensor {
	generatedMetrics := make(map[uint64]commons.Measurement)

	return &Sensor{
		ID:               cfg.SensorID,
		natsConfig:       cfg.Nats,
		sensorConfig:     newSensorConfig(),
		natsClient:       natsClient,
		generatedMetrics: generatedMetrics,
	}
}

func (s *Sensor) generateFakeMeasurement() commons.Measurement {
	return commons.Measurement{
		LectureID: uint64(time.Now().UnixNano()),
		SensorID:  s.ID,
		Lecture:   fmt.Sprintf("Fake Lecture %d", rand.Intn(1000)),
	}
}

func (s *Sensor) postMeasurement(measurement commons.Measurement) error {
	s.generatedMetrics[measurement.LectureID] = measurement

	return s.natsClient.PostMeasurement(commons.MeasurementToBytes(measurement))
}
