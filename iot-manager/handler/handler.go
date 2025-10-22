// Package handler exposes an HTTP REST API for managing IoT devices.
package handler

import "github.com/alexrodfe/iot-poc/iot-manager/adapter/clients"

type Handler struct {
	natsManager clients.NatsClient
}

func New(natsManager clients.NatsClient) *Handler {
	return &Handler{natsManager: natsManager}
}

func (h *Handler) StartService() error {

	return nil
}
