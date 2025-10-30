package main

import (
	"log"

	"github.com/alexrodfe/iot-poc/iot-manager/adapter/implementations"
	"github.com/alexrodfe/iot-poc/iot-manager/config"
	"github.com/alexrodfe/iot-poc/iot-manager/handler"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	natsClient, err := implementations.NewNatsManager(cfg.Nats)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.New(natsClient)

	handler.StartService()
}
