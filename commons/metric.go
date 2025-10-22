// Package commons declares public data types shared between the devices and the IoT manager.
package commons

type Metric struct {
	ID          string
	Temperature float64
	Humidity    float64
	Pressure    float64
	Timestamp   int64
}
