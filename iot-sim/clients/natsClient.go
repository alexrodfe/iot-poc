// Package clients implements technology clients used in the IoT sim
package clients

import (
	"errors"
	"fmt"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sim/config"
	"github.com/nats-io/nats.go"
)

type HandleCommandFunc func(command commons.SensorCommand) error

type Nats interface {
	SubscribeToMeasurements(sensorID string) error
	EditConfig(sensorID string, command commons.SensorCommand) error
	EditMeasurement(sensorID string, command commons.SensorCommand) error
	RetrieveSensorConfig(sensorID string) (string, error)
}

var _ Nats = (*natsClient)(nil)

type natsClient struct {
	cfg config.NatsConfig

	js nats.JetStreamContext
}

func NewNatsClient(cfg config.NatsConfig) (Nats, error) {
	natsConn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	js, err := natsConn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS JetStream: %w", err)
	}

	return &natsClient{cfg: cfg, js: js}, nil
}

func (n *natsClient) SubscribeToMeasurements(sensorID string) error {
	_, err := n.js.Subscribe(commons.GetMeasurementsStream(sensorID), func(msg *nats.Msg) {
		measurement, err := commons.BytesToMeasurement(msg.Data)
		if err != nil {
			fmt.Printf("error unmarshaling measurement: %v\n", err)
			msg.Ack()
			return
		}
		fmt.Printf("Received measurement: %s\n", measurement.String())
		msg.Ack()
	}, nats.Durable("measurements-durable"), nats.ManualAck())
	if err != nil {
		return fmt.Errorf("error subscribing to measurements: %w", err)
	}
	return nil
}

func (n *natsClient) EditConfig(sensorID string, command commons.SensorCommand) error {
	data := commons.SensorCommandToBytes(command)

	_, err := n.js.Publish(commons.GetSensorStream(sensorID), data)
	if err != nil {
		return fmt.Errorf("error publishing command: %w", err)
	}

	return nil
}

func (n *natsClient) EditMeasurement(sensorID string, command commons.SensorCommand) error {
	data := commons.SensorCommandToBytes(command)

	_, err := n.js.Publish(commons.GetSensorStream(sensorID), data)
	if err != nil {
		return fmt.Errorf("error publishing command: %w", err)
	}

	return nil
}

func (n *natsClient) RetrieveSensorConfig(sensorID string) (string, error) {
	kv, err := n.js.KeyValue(commons.GetSensorStream(sensorID))
	if err != nil {
		return "", fmt.Errorf("kv bucket %q not available: %w", commons.GetSensorStream(sensorID), err)
	}

	entry, err := kv.Get("config")
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("error retrieving config: %w", err)
	}

	return string(entry.Value()), nil
}
