package sensor

import (
	"fmt"

	"github.com/alexrodfe/iot-poc/commons"
)

func (s *Sensor) StartHandler() error {
	err := s.natsClient.StartHandler(s.handleCommand)
	if err != nil {
		return fmt.Errorf("error starting NATS handler: %w", err)
	}

	return nil
}

func (s *Sensor) handleCommand(command commons.SensorCommand) error {
	switch command.CommandID {
	case commons.CommandUpdateConfig:
		s.handleUpdateConfig(command)
	case commons.CommandUpdateMetric:
		s.handleUpdateMetric(command)
	default:
		return fmt.Errorf("unknown command: %d", command.CommandID)
	}
	return nil
}

func (s *Sensor) handleUpdateConfig(command commons.SensorCommand) {
	s.SetInterval(command.NewIntervalMilliseconds)
}

func (s *Sensor) handleUpdateMetric(command commons.SensorCommand) {
	s.generatedMetrics[command.MetricID] = commons.Measurement{
		LectureID: command.MetricID,
		SensorID:  s.ID,
		Lecture:   command.NewMeasurementValue,
	}
}
