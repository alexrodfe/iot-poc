package implementations

import (
	"fmt"

	"github.com/alexrodfe/iot-poc/iot-manager/adapter/clients"
	"github.com/nats-io/nats.go"
)

var _ clients.NatsClient = (*NatsManager)(nil)

type NatsManager struct {
	js nats.JetStreamContext
}

func NewNatsManager(js nats.JetStreamContext) *NatsManager {
	return &NatsManager{js: js}
}

func InitJetStreamConnection(url string) (nats.JetStreamContext, error) {
	natsConn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %w", err)
	}

	js, err := natsConn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS JetStream: %w", err)
	}

	return js, nil
}
