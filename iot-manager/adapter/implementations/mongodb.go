// Package implementations contains the implementations of the adapters for external technologies.
package implementations

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBClient struct {
	client *mongo.Client
}
