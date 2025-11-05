package sensor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alexrodfe/iot-poc/commons"
)

func (s *Sensor) generateFakeMeasurement() commons.Measurement {
	return commons.Measurement{
		LectureID: uint64(time.Now().UnixNano()),
		SensorID:  s.ID,
		Lecture:   fmt.Sprintf("Fake Lecture %d", rand.Intn(1000)),
	}
}

func (s *Sensor) postMeasurement(measurement commons.Measurement) error {
	s.generatedMeasurements[measurement.LectureID] = measurement

	return s.natsClient.PostMeasurement(commons.MeasurementToBytes(measurement))
}

func (s *Sensor) editMeasurement(metricID uint64, newValue string) {
	s.generatedMeasurements[metricID] = commons.Measurement{
		LectureID: metricID,
		SensorID:  s.ID,
		Lecture:   newValue,
	}
}
