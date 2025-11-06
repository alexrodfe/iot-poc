package commons

type SensorEnv struct {
	SensorID string
	NatsEnv  SensorNatsEnv
}

type SensorNatsEnv struct {
	SensorStream       string
	MeasurementsStream string
}

func NewNatsEnv(sensorID string) SensorNatsEnv {
	return SensorNatsEnv{
		SensorStream:       GetSensorStream(sensorID),
		MeasurementsStream: GetMeasurementsStream(sensorID),
	}
}

func GetSensorStream(sensorID string) string {
	return sensorID + "_sensor_stream"
}

func GetMeasurementsStream(sensorID string) string {
	return sensorID + "_measurements_stream"
}
