package sensor

import (
	"context"
	"log"
	"time"
)

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

			log.Printf("✅ Published measurement id=%d sensor=%d", m.LectureID, m.SensorID)
		}
	}
}
