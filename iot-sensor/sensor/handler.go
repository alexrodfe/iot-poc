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
	case commons.CommandUpdateMeasurement:
		s.handleUpdateMetric(command)
	default:
		return fmt.Errorf("unknown command: %d", command.CommandID)
	}
	return nil
}

func (s *Sensor) handleUpdateConfig(command commons.SensorCommand) {
	fmt.Printf("-----Change of interval detected, new interval set to %d milliseconds-----\n", command.NewIntervalMilliseconds)
	newSensorConfig := s.SetInterval(command.NewIntervalMilliseconds)
	s.natsClient.UpdateSensorConfig(newSensorConfig.ToString())
}

func (s *Sensor) handleUpdateMetric(command commons.SensorCommand) {
	fmt.Printf("-----Change of metric detected, new lecture set to %s for LectureID %d-----\n", command.NewMeasurementValue, command.MeasurementID)
	s.editMeasurement(command.MeasurementID, command.NewMeasurementValue)
}
