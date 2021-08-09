package config

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbSession *mongo.Client
var Database string

func SetDbSession(session *mongo.Client, databaseName string) {

	Database = databaseName
	DbSession = session
}

func SetUpDatabase() {

	var mongoURI string
	mongoURI = "mongodb+srv://aqua:aqua@covid-19.zkr0h.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI).
		SetConnectTimeout(2 * time.Second).
		SetMaxPoolSize(500).
		SetSocketTimeout(2 * time.Second)

	// Connect to MongoDB
	session, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatalf("connection with atlas mongo failed with error : %v", err.Error())
	}

	// Check the connection
	err = session.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatalf("connection with atlas mongo failed with error : %v", err.Error())
	}
	SetDbSession(session, "covid19")
}

func GetDatabase() *mongo.Database {

	return DbSession.Database(Database)
}
