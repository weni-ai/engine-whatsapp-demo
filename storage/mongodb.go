package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbName = "whatsapp-router"

func NewDB() *mongo.Database {
	//TODO move to env config
	uri := "mongodb://admin:admin@127.0.0.1:27017/?appName=whatsapp-router"
	options := options.Client().ApplyURI(uri)
	connection, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		panic(err)
	}

	if err := connection.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected and pinged on mongodb.")
	}

	db := connection.Database(dbName)
	return db
}

func CloseDB(db *mongo.Database) {
	if err := db.Client().Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on close MongoDB: %v", err)
	}
}
