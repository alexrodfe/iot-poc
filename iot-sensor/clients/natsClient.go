// Package clients implements technology clients used in the IoT sensor
package clients

import (
	"fmt"

	"github.com/alexrodfe/iot-poc/iot-sensor/config"
	"github.com/nats-io/nats.go"
)

type Nats interface {
	PostMeasurement(data []byte) error
}

var _ Nats = (*natsClient)(nil)

type natsClient struct {
	cfg config.NatsConfig

	js nats.JetStreamContext
}

func NewNatsClient(cfg config.NatsConfig) (*natsClient, error) {
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
		return nil
	}

	// Create the sensor stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.SensorStream,
		Subjects: []string{"sensor.>"},
	})
	if err != nil {
		return fmt.Errorf("error creating sensor stream: %w", err)
	}

	return nil
}

func createMeasurementsStream(js nats.JetStreamContext, cfg config.NatsConfig) error {
	// If stream already exists, do nothing
	stream, err := js.StreamInfo(cfg.MeasurementsStream)
	if err == nil && stream != nil {
		return nil
	}

	// Create the measurements stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.MeasurementsStream,
		Subjects: []string{"measurements.>"},
	})
	if err != nil {
		return fmt.Errorf("error creating measurements stream: %w", err)
	}

	return nil
}

func (n *natsClient) PostMeasurement(data []byte) error {
	_, err := n.js.Publish(n.cfg.MeasurementsStream, data)
	if err != nil {
		return fmt.Errorf("error publishing measurement to NATS: %w", err)
	}
	return nil
}
