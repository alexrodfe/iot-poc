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

	natsClient, err := clients.NewNatsClient(cfg.Nats)
	if err != nil {
		panic(err)
	}

	cmd := os.Args[1]

	switch cmd {

	case "subscribe":
		subCmd := flag.NewFlagSet("subscribe", flag.ExitOnError)
		sensorID := subCmd.String("sensorID", "", "sensor identifier (required)")
		_ = subCmd.Parse(os.Args[2:])
		requireSensorID(subCmd, *sensorID)

		fmt.Printf("Subscribing for sensorID=%s\n", *sensorID)
		natsClient.SubscribeToMeasurements(*sensorID)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Shutdown signal received. Exiting.")

	case "edit_config":
		editCmd := flag.NewFlagSet("edit_config", flag.ExitOnError)
		sensorID := editCmd.String("sensorID", "", "sensor identifier (required)")
		interval := editCmd.Int("interval", 0, "new interval (seconds)")
		_ = editCmd.Parse(os.Args[2:])
		requireSensorID(editCmd, *sensorID)

		command := commons.NewUpdateConfigCommand(uint64(*interval))
		fmt.Printf("Editing config sensorID=%s interval=%d\n", *sensorID, *interval)
		natsClient.EditConfig(*sensorID, command)

	case "edit_measurement":
		editCmd := flag.NewFlagSet("edit_measurement", flag.ExitOnError)
		sensorID := editCmd.String("sensorID", "", "sensor identifier (required)")
		metricID := editCmd.Uint64("id", 0, "metric ID")
		newValue := editCmd.String("value", "", "new measurement value")
		_ = editCmd.Parse(os.Args[2:])
		requireSensorID(editCmd, *sensorID)

		command := commons.NewUpdateMeasurementCommand(*metricID, *newValue)
		fmt.Printf("Editing measurement sensorID=%s metricID=%d value=%s\n", *sensorID, *metricID, *newValue)
		natsClient.EditMeasurement(*sensorID, command)

	case "get_config":
		getCmd := flag.NewFlagSet("get_config", flag.ExitOnError)
		sensorID := getCmd.String("sensorID", "", "sensor identifier (required)")
		_ = getCmd.Parse(os.Args[2:])
		requireSensorID(getCmd, *sensorID)

		cfgStr, err := natsClient.RetrieveSensorConfig(*sensorID)
		if err != nil {
			fmt.Println("Error retrieving sensor config:", err)
			return
		}
		if cfgStr == "" {
			fmt.Printf("No sensor config found for sensorID=%s\n", *sensorID)
			return
		}
		fmt.Printf("Sensor config for %s: %s\n", *sensorID, cfgStr)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func requireSensorID(fs *flag.FlagSet, sensorID string) {
	if sensorID == "" {
		fmt.Println("Error: --sensorID is required")
		fs.Usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  iot-sim subscribe --sensorID <id>")
	fmt.Println("  iot-sim get_config --sensorID <id>")
	fmt.Println("  iot-sim edit_config --sensorID <id> --interval <seconds>")
	fmt.Println("  iot-sim edit_measurement --sensorID <id> --id <measurement_id> --value <new_value>")
}
