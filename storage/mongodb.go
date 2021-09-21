package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/weni/whatsapp-router/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbName = "whatsapp-router"

func NewDB() *mongo.Database {
	//TODO move to env config
	uri := "mongodb://admin:admin@127.0.0.1:27017/?appName=whatsapp-router"
	options := options.Client().ApplyURI(uri)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	connection, err := mongo.Connect(ctx, options)
	if err != nil {
		logger.Error("mongodb FAIL")
		panic(err.Error())
	}

	if err := connection.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Error("mongodb FAIL")
		panic(err.Error())
	} else {
		logger.Info("mongodb OK")
	}

	db := connection.Database(dbName)
	return db
}

func CloseDB(db *mongo.Database) {
	if err := db.Client().Disconnect(context.TODO()); err != nil {
		logger.Error(fmt.Sprintf("Error on close MongoDB: %v", err))
		log.Fatal()

	}
}
