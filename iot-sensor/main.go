// package main will load up all dependencies and run the IoT sensor
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexrodfe/iot-poc/iot-sensor/clients"
	"github.com/alexrodfe/iot-poc/iot-sensor/config"
	sensPck "github.com/alexrodfe/iot-poc/iot-sensor/sensor"
)

func main() {
	// Build dependencies
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	nc, err := clients.NewNatsClient(cfg.Nats)
	if err != nil {
		panic(err)
	}

	sensor := sensPck.New(cfg, nc)

	// Prepare context that cancels on SIGTERM / SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signals
		cancel()
	}()

	if err := sensor.StartHandler(); err != nil {
		log.Fatalf("Fatal error starting sensor handler: %v", err)
	}

	if err := sensor.Run(ctx); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
