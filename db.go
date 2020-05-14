package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var Context = context.TODO()
var Client *mongo.Client
var Database *mongo.Database
var AdherentsCollection *mongo.Collection

func Connect() {
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		uri = "mongodb://vps.ribes.ovh:27017"
	}
	Client, _ := mongo.Connect(Context, options.Client().ApplyURI(uri))
	handleError(Client.Ping(Context, nil))

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "ami"
	}
	Database = Client.Database(databaseName)
	AdherentsCollection = Database.Collection("adherents")
}
func Disconnect() {
	handleError(Client.Disconnect(Context))
}
