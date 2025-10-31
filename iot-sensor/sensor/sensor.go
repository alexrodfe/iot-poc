// Package sensor declares the IoT sensor
package sensor

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sensor/clients"
	"github.com/alexrodfe/iot-poc/iot-sensor/config"
)

type Sensor struct {
	ID uint64

	sensorConfig sensorConfig
	natsConfig   config.NatsConfig

	natsClient clients.Nats
}

func New(id uint64, natsClient clients.Nats) *Sensor {
	return &Sensor{
		ID:           id,
		sensorConfig: newSensorConfig(),
		natsClient:   natsClient,
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
	return s.natsClient.PostMeasurement(commons.MeasurementToBytes(measurement))
}

// Run starts the publishing loop.
func (s *Sensor) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.GetInterval())
	defer ticker.Stop()

	log.Printf("🚀 Sensor started. ")

	for {
		currentInterval := s.GetInterval()

		select {
		case <-ctx.Done():
			log.Println("🛑 Shutdown signal received, stopping sensor loop.")
			return nil
		case <-time.After(currentInterval):
			m := s.generateFakeMeasurement()

			err := s.postMeasurement(m)
			if err != nil {
				log.Println("⚠️ Failed to post measurement, skipping.")
				continue
			}

			log.Printf("✅ Published measurement id=%d sensor=%d", m.LectureID, m.SensorID)
		}
	}
}
