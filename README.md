# iot-poc

Proof of concept for an Internet of Things platform. It simulates IoT sensors producing measurements and an admin/simulator component that sends commands. Communication is via NATS JetStream. Multiple Go modules are managed with a `go.work` workspace:

- `commons`: shared types (Measurement, SensorCommand).
- `iot-sensor`: fake sensor service (publishes measurements, consumes commands, writes config to KV).
- `iot-sim`: admin / simulator (edits sensor config, can subscribe to measurements, fetch sensor config).

## Prerequisites

- Go ≥ 1.22 (workspace support).
- Docker (optional, for running NATS).
- NATS Server (JetStream enabled).

## Start NATS

Use the helper script (JetStream enabled):

```sh
bash hack/run-nats.sh
```

Or manually:

```sh
docker run --rm -p 4222:4222 -p 8222:8222 nats:latest -js -m 8222
```

## Environment

Create `.env`:

```sh
NATS_URL=nats://localhost:4222
```

You can add/customize streams & subjects if your config package supports them:

```sh
NATS_SENSOR_STREAM=SENSOR_STREAM
NATS_SENSOR_SUBJECT=sensor.commands
NATS_MEASUREMENTS_STREAM=MEASUREMENTS_STREAM
NATS_MEASUREMENTS_SUBJECT=sensor.measurements
```

## Workspace setup

```sh
go work init ./commons ./iot-sensor ./iot-sim
go work sync
```

## Run a sensor

Each sensor requires a unique ID:

```sh
cd iot-sensor
go run . --sensorID sensor-1
```

The sensor:

- Publishes measurements to the measurements stream.
- Consumes commands from the command stream (ack deletes due to WorkQueue retention).
- Updates its config in a JetStream KV bucket.

## Run simulator commands

```sh
cd iot-sim
# Subscribe to measurements (blocks until Ctrl+C)
go run . subscribe --sensorID sensor-1

# Get current sensor config
go run . get_config --sensorID sensor-1

# Edit sensor config (interval example)
go run . edit_config --sensorID sensor-1 --interval 5

# Edit a measurement (example placeholder)
go run . edit_measurement --sensorID sensor-1 --id 123 --value "new-value"
```

## Architecture (summary)

- `iot-sensor` publishes JSON `Measurement` messages (JetStream stream: `MEASUREMENTS_STREAM`).
- `iot-sim` publishes JSON `SensorCommand` messages (JetStream stream: `SENSOR_STREAM`).
- Sensor writes config to KV bucket (`SENSOR_STREAM`, key: `config`), simulator reads it.
- Streams:
  - Measurements: retention = limits (history kept until size/age limits).
  - Commands: retention = work queue (message removed after explicit ack).

## Testing

Integration tests (JetStream + Testcontainers) example (run from repo root):

```sh
go test -tags=integration ./...
```

## Directory layout

```sh
commons/          Shared domain types
iot-sensor/       Sensor executable + NATS client + dynamic loop
iot-sim/          Simulator / admin CLI
hack/             Utility scripts (run-nats.sh)
.env              Environment variables
```

## License

MIT (see LICENSE).

## Next steps

- Persist measurements to a database.
- Add validation / schema versioning.
- Expose a REST API for management in place of CLI.
