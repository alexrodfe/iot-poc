// Package commons declares public data types shared between the sensors and the IoT manager.
package commons

import (
	"encoding/json"
	"fmt"
)

type Measurement struct {
	LectureID uint64 `json:"lecture_id"`
	SensorID  string `json:"sensor_id"`
	Lecture   string `json:"lecture"`
}

func (m Measurement) String() string {
	return fmt.Sprintf("Measurement{LectureID:%d SensorID:%s Lecture:%q}", m.LectureID, m.SensorID, m.Lecture)
}

func MeasurementToBytes(m Measurement) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return b
}

func BytesToMeasurement(data []byte) (Measurement, error) {
	var m Measurement
	err := json.Unmarshal(data, &m)
	return m, err
}
