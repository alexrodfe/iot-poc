package sensor

import (
	"fmt"
	"time"
)

type sensorConfig struct {
	millisecondsInterval uint64
}

// newSensorConfig creates a new sensorConfig with default values.
func newSensorConfig() sensorConfig {
	return sensorConfig{
		millisecondsInterval: 10000,
	}
}

func (s *Sensor) GetInterval() time.Duration {
	return time.Millisecond * time.Duration(s.sensorConfig.millisecondsInterval)
}

func (s *Sensor) SetInterval(newInterval uint64) sensorConfig {
	s.sensorConfig.millisecondsInterval = newInterval
	return s.sensorConfig
}

func (sc *sensorConfig) ToString() string {
	return fmt.Sprintf("SensorConfig{millisecondsInterval: %d}", sc.millisecondsInterval)
}
