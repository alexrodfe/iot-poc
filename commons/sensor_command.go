package commons

import "encoding/json"

type Commands uint8

const (
	CommandUpdateConfig Commands = 1
	CommandUpdateMetric Commands = 2
)

type SensorCommand struct {
	CommandID Commands `json:"command_id"`

	MetricID            uint64 `json:"metric_id"`
	NewMeasurementValue string `json:"new_measurement_value"`

	NewIntervalMilliseconds uint64 `json:"new_interval_milliseconds"`
}

func NewUpdateConfigCommand(newInterval uint64) SensorCommand {
	return SensorCommand{
		CommandID:               CommandUpdateConfig,
		NewIntervalMilliseconds: newInterval,
	}
}

func NewUpdateMetricCommand(metricID uint64, newValue string) SensorCommand {
	return SensorCommand{
		CommandID:           CommandUpdateMetric,
		MetricID:            metricID,
		NewMeasurementValue: newValue,
	}
}

func SensorCommandToBytes(sc SensorCommand) []byte {
	b, err := json.Marshal(sc)
	if err != nil {
		return []byte{}
	}
	return b
}

func BytesToSensorCommand(data []byte) (SensorCommand, error) {
	var m SensorCommand
	err := json.Unmarshal(data, &m)
	return m, err
}
