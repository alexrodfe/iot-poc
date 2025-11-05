// Package clients implements technology clients used in the IoT sensor
package clients

import (
	"encoding/json"
	"fmt"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sensor/config"
	"github.com/nats-io/nats.go"
)

type HandleCommandFunc func(command commons.SensorCommand) error

type Nats interface {
	PostMeasurement(data []byte) error
	StartHandler(handleCommand HandleCommandFunc) error
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

	err = createSensorStream(js, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating sensor stream: %w", err)
	}

	err = createMeasurementsStream(js, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating measurements stream: %w", err)
	}

	return &natsClient{cfg: cfg, js: js}, nil
}

func createSensorStream(js nats.JetStreamContext, cfg config.NatsConfig) error {
	// If stream already exists, do nothing
	stream, err := js.StreamInfo(cfg.SensorStream)
	if err == nil && stream != nil {
		fmt.Printf("Sensor stream %s already exists\n", cfg.SensorStream)
		return nil
	}

	// Create the sensor stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      cfg.SensorStream,
		Subjects:  []string{cfg.SensorSubject},
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		return fmt.Errorf("error creating sensor stream: %w", err)
	}

	fmt.Printf("Sensor stream %s succesfuly created\n", cfg.SensorStream)

	return nil
}

func createMeasurementsStream(js nats.JetStreamContext, cfg config.NatsConfig) error {
	// If stream already exists, do nothing
	stream, err := js.StreamInfo(cfg.MeasurementsStream)
	if err == nil && stream != nil {
		fmt.Printf("Measurements stream %s already exists\n", cfg.MeasurementsStream)
		return nil
	}

	// Create the measurements stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.MeasurementsStream,
		Subjects: []string{cfg.MeasurementsSubject},
	})
	if err != nil {
		return fmt.Errorf("error creating measurements stream: %w", err)
	}

	fmt.Printf("Measurements stream %s succesfuly created\n", cfg.MeasurementsStream)

	return nil
}

func (n *natsClient) PostMeasurement(data []byte) error {
	_, err := n.js.Publish(n.cfg.MeasurementsSubject, data)
	if err != nil {
		return fmt.Errorf("error publishing measurement to NATS: %w", err)
	}
	return nil
}

func (n *natsClient) StartHandler(handleCommand HandleCommandFunc) error {
	_, err := n.js.Subscribe(n.cfg.SensorSubject, func(m *nats.Msg) {
		var command commons.SensorCommand
		if err := json.Unmarshal(m.Data, &command); err != nil {
			fmt.Println("Warning, unable to unmarshall incoming message, ignoring..")
			m.Ack()
			return
		}

		if err := handleCommand(command); err != nil {
			fmt.Printf("Error handling incoming command: %v\n", err)
			m.Ack()
			return
		}

		m.Ack()
	})

	if err != nil {
		return fmt.Errorf("error subscribing to sensor stream: %w", err)
	}

	return nil
}
