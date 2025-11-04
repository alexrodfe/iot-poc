// Package clients implements technology clients used in the IoT sim
package clients

import (
	"fmt"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sim/config"
	"github.com/nats-io/nats.go"
)

type HandleCommandFunc func(command commons.SensorCommand) error

type Nats interface {
	SubscribeToMeasurements() error
	EditConfig(commons.SensorCommand) error
	EditMeasurement(commons.SensorCommand) error
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

func (n *natsClient) SubscribeToMeasurements() error {
	_, err := n.js.Subscribe(n.cfg.MeasurementsSubject, func(msg *nats.Msg) {
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

func (n *natsClient) EditConfig(command commons.SensorCommand) error {
	data := commons.SensorCommandToBytes(command)

	_, err := n.js.Publish(n.cfg.SensorSubject, data)
	if err != nil {
		return fmt.Errorf("error publishing command: %w", err)
	}

	return nil
}

func (n *natsClient) EditMeasurement(command commons.SensorCommand) error {
	data := commons.SensorCommandToBytes(command)

	_, err := n.js.Publish(n.cfg.MeasurementsSubject, data)
	if err != nil {
		return fmt.Errorf("error publishing command: %w", err)
	}

	return nil
}
