//go:build integration

package clients_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alexrodfe/iot-poc/commons"
	"github.com/alexrodfe/iot-poc/iot-sensor/clients"
	"github.com/alexrodfe/iot-poc/iot-sensor/config"
)

const (
	sensorID = "test-sensor"
)

// Esta es la unica test suite que hay en el proyecto por ahora, pero la idea es que sirva de ejemplo para replicarla en el resto de componentes.
// Este ejemplo me permite demostrar cómo queda una bateria de tests para un paquete concreto bajo una suite de tests.

type NatsClientTestSuite struct {
	suite.Suite

	natsContainer testcontainers.Container
	js            nats.JetStreamContext

	natsClient clients.Nats
	natsConfig config.NatsConfig
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(NatsClientTestSuite))
}

func (s *NatsClientTestSuite) SetupSuite() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "nats:2.8.1",
		Cmd:          []string{"-js"},
		ExposedPorts: []string{"4222/tcp", "8222/tcp"},
		WaitingFor: wait.ForAll(wait.ForLog("Server is ready")).
			WithDeadline(time.Minute * 3),
	}

	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	natsEndpoint, err := natsContainer.Endpoint(ctx, "")
	s.Require().NoError(err)

	natsCfg := config.NewNatsConfig(sensorID, natsEndpoint)
	natsClient, err := clients.NewNatsClient(natsCfg)
	s.Require().NoError(err)

	natsConn, err := nats.Connect(natsEndpoint)
	s.Require().NoError(err)

	js, err := natsConn.JetStream()
	s.Require().NoError(err)

	s.natsContainer = natsContainer
	s.js = js

	s.natsClient = natsClient
	s.natsConfig = natsCfg
}

func (s *NatsClientTestSuite) TearDownSuite() {
	if err := s.natsContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err.Error())
	}
}

func (s *NatsClientTestSuite) TestStreamsExistOnStartup() {
	// GIVEN suite setup

	// THEN
	sensorStreamInfo, err := s.js.StreamInfo(s.natsConfig.SensorStream)
	s.Require().NoError(err)
	s.Assert().NotNil(sensorStreamInfo)

	measurementsStreamInfo, err := s.js.StreamInfo(s.natsConfig.MeasurementsStream)
	s.Require().NoError(err)
	s.Assert().NotNil(measurementsStreamInfo)
}

func (s *NatsClientTestSuite) TestPostMeasurementPublishes() {
	// GIVEN
	data, _ := json.Marshal(commons.Measurement{
		LectureID: 123,
		SensorID:  sensorID,
		Lecture:   "dummy",
	})

	// WHEN
	err := s.natsClient.PostMeasurement(data)
	s.Require().NoError(err)

	// THEN
	info, err := s.js.StreamInfo(s.natsConfig.MeasurementsStream)
	s.Require().NoError(err)
	s.Assert().Equal(uint64(1), info.State.Msgs)
}

func (s *NatsClientTestSuite) TestStartHandlerAckDeletes() {
	// GIVEN
	err := s.natsClient.StartHandler(func(cmd commons.SensorCommand) error {
		return nil
	})
	s.Require().NoError(err)

	// WHEN
	cmd := commons.SensorCommand{CommandID: commons.CommandUpdateConfig, NewIntervalMilliseconds: 1000}
	raw, _ := json.Marshal(cmd)
	_, err = s.js.Publish(s.natsConfig.SensorStream, raw)
	s.Require().NoError(err)

	time.Sleep(500 * time.Millisecond)

	// THEN
	info, err := s.js.StreamInfo(s.natsConfig.SensorStream)
	s.Require().NoError(err)
	s.Assert().Empty(info.State.Msgs)
}

func (s *NatsClientTestSuite) TestUpdateSensorConfigStoresKV() {
	// WHEN
	err := s.natsClient.UpdateSensorConfig(`{"interval_ms":2500}`)
	s.Require().NoError(err)

	// THEN
	kv, err := s.js.KeyValue(s.natsConfig.SensorStream)
	s.Require().NoError(err)

	entry, err := kv.Get("config")
	s.Require().NoError(err)
	s.Assert().Contains(string(entry.Value()), `"interval_ms":2500`)
}
