package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sim/clients"
	"github.com/alexrodfe/iot-poc/iot-sim/config"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Initialize NATS client
	natsClient, err := clients.NewNatsClient(cfg.Nats)
	if err != nil {
		panic(err)
	}

	cmd := os.Args[1]

	switch cmd {
	case "subscribe":
		natsClient.SubscribeToMeasurements()

		// Wait for termination signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Shutdown signal received. Exiting.")

	case "edit_config":
		editCmd := flag.NewFlagSet("edit_config", flag.ExitOnError)
		interval := editCmd.Int("interval", 0, "new interval (seconds)")
		_ = editCmd.Parse(os.Args[2:])

		command := commons.NewUpdateConfigCommand(uint64(*interval))

		natsClient.EditConfig(command)

	case "edit_measurement":
		editCmd := flag.NewFlagSet("edit_measurement", flag.ExitOnError)
		metricID := editCmd.Uint64("id", 0, "metric ID")
		newValue := editCmd.String("value", "", "new measurement value")
		_ = editCmd.Parse(os.Args[2:])

		command := commons.NewUpdateMetricCommand(*metricID, *newValue)

		natsClient.EditMeasurement(command)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("iot-sim subscribe")
	fmt.Println("iot-sim edit_config --interval <new_interval>")
	fmt.Println("iot-sim edit_measurement --id <sensor_id> --value <new_value>")
}
