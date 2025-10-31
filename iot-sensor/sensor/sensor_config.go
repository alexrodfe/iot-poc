package sensor

import "time"

type sensorConfig struct {
	millisecondsInterval int
}

// newsensorConfig creates a new sensorConfig with default values.
func newSensorConfig() sensorConfig {
	return sensorConfig{
		millisecondsInterval: 1000,
	}
}

func (s *Sensor) GetInterval() time.Duration {
	return time.Millisecond * time.Duration(s.sensorConfig.millisecondsInterval)
}
