package main

import (
	"log"

	"github.com/alexrodfe/iot-poc/iot-manager/adapter/implementations"
	"github.com/alexrodfe/iot-poc/iot-manager/handler"
	"github.com/spf13/viper"
)

func main() {
	// TODO: load nats url through viper from config file or env variable
	js, err := implementations.InitJetStreamConnection(viper.GetString("nats url goes here"))
	if err != nil {
		log.Fatal(err)
	}

	natsClient := implementations.NewNatsManager(js)
	handler := handler.New(natsClient)

	handler.StartService()
}
