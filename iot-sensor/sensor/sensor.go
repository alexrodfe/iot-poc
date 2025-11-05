// Package sensor declares the IoT sensor
package sensor

import (
	"context"
	"fmt"
	"log"
	"os"
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

	generatedMeasurements map[uint64]commons.Measurement
}

func New(cfg *config.Config, natsClient clients.Nats) *Sensor {
	generatedMetrics := make(map[uint64]commons.Measurement)

	return &Sensor{
		ID:                    cfg.SensorID,
		natsConfig:            cfg.Nats,
		sensorConfig:          newSensorConfig(),
		natsClient:            natsClient,
		generatedMeasurements: generatedMetrics,
	}
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
				log.Println(err)
				continue
			}

			log.Printf("✅ Published measurement id=%d sensor=%s", m.LectureID, m.SensorID)
		}
	}
}

func (s *Sensor) HandleShutdown() error {
	log.Println("Handling sensor shutdown...")

	// output all generated metrics into a file
	if len(s.generatedMeasurements) == 0 {
		log.Println("No generated metrics to persist.")
		return nil
	}

	filename := fmt.Sprintf("generated_metrics_%d.txt", time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create metrics file: %w", err)
	}
	defer f.Close()

	for k, m := range s.generatedMeasurements {
		if _, err := fmt.Fprintf(f, "%d\t%s\t%s\n", k, m.SensorID, m.Lecture); err != nil {
			return fmt.Errorf("error writing metric %d: %w", k, err)
		}
	}

	log.Printf("Persisted %d metrics to %s\n", len(s.generatedMeasurements), filename)
	return nil
}
