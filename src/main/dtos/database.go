package dtos

import "go.mongodb.org/mongo-driver/mongo"

type DbConnections struct {
	PrimaryConnection   *mongo.Database
	SecondaryConnection *mongo.Database
}
