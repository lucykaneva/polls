package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var Client mongo.Client

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	local, err := mongo.Connect(options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/mongo?authSource=admin"))
	if err != nil {
		return err
	}
	Client = *local

	if err = Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB!")
	return nil
}
