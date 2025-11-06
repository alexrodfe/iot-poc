// package main will load up all dependencies and run the IoT sensor
package main

import (
	"context"
	"flag"
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
	sensorID := flag.String("sensorID", "", "sensor identifier (required)")
	flag.Parse()

	if *sensorID == "" {
		fmt.Println("Error: --sensorID is required")
		usage()
		os.Exit(2)
	}

	// Build dependencies
	cfg, err := config.New(*sensorID)
	if err != nil {
		panic(err)
	}

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

	if err := sensor.HandleShutdown(); err != nil {
		log.Fatalf("Fatal error handling shutdown: %v", err)
	}

	fmt.Println("Sensor shut down gracefully")
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("iot-sensor --sensorID <id>")
}
