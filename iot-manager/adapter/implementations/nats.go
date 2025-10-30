package implementations

import (
	"fmt"

	"github.com/alexrodfe/iot-poc/iot-manager/adapter/clients"
	"github.com/alexrodfe/iot-poc/iot-manager/config"
	"github.com/nats-io/nats.go"
)

var _ clients.NatsClient = (*NatsManager)(nil)

type NatsManager struct {
	js nats.JetStreamContext
}

func NewNatsManager(cfg config.NatsConfig) (*NatsManager, error) {
	natsConn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	js, err := natsConn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS JetStream: %w", err)
	}

	err = createAdminStream(js, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating admin stream: %w", err)
	}

	return &NatsManager{js: js}, nil
}

func createAdminStream(js nats.JetStreamContext, cfg config.NatsConfig) error {
	// If stream already exists, do nothing
	stream, err := js.StreamInfo(cfg.AdminStream)
	if err == nil && stream != nil {
		return nil
	}

	// Create the admin stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.AdminStream,
		Subjects: []string{"admin.>"},
	})
	if err != nil {
		return fmt.Errorf("error creating admin stream: %w", err)
	}

	return nil
}
