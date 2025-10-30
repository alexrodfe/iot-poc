// Package handler exposes an HTTP REST API for managing IoT devices.
package handler

import "github.com/alexrodfe/iot-poc/iot-manager/adapter/clients"

type NatsHandler struct {
	natsManager clients.NatsClient
}

func New(natsManager clients.NatsClient) *NatsHandler {
	return &NatsHandler{natsManager: natsManager}
}

func (h *NatsHandler) StartService() error {

	return nil
}
