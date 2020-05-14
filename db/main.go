package db

import (
	"context"
	"encoding/gob"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	log.Printf("Connexion à to %s", uri)
	Client, _ := mongo.Connect(Context, options.Client().ApplyURI(uri))
	handleError(Client.Ping(Context, nil))

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "ami"
	}
	Database = Client.Database(databaseName)
	AdherentsCollection = Database.Collection("adherents")
	log.Print("connecté à Mongo")

	gob.Register(Adherent{})
}
func Disconnect() {
	handleError(Client.Disconnect(Context))
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
