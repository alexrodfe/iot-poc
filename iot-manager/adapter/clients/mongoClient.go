// Package clients contains all interfaces for the usage of external technologies.
package clients

type MongoClient interface {
	Connect() error
	Close() error

	GetDeviceEntriesByID(id string) (string, error)
}
