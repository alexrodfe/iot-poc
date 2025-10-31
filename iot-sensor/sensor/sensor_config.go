package sensor

import (
	"time"
)

type sensorConfig struct {
	millisecondsInterval uint64
}

// newSensorConfig creates a new sensorConfig with default values.
func newSensorConfig() sensorConfig {
	return sensorConfig{
		millisecondsInterval: 1000,
	}
}

func (s *Sensor) GetInterval() time.Duration {
	return time.Millisecond * time.Duration(s.sensorConfig.millisecondsInterval)
}

func (s *Sensor) SetInterval(newInterval uint64) {
	s.sensorConfig.millisecondsInterval = newInterval
}
