package main

import (
	"context"
	"fmt"
	"log"
	"polls/src/pckg/db"
	handlers "polls/src/pckg/handlers"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var Client *mongo.Client

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}
	// collection := client.Database("pollsdb").Collection("polls")
	// models.InitCollection(collection)

	router := gin.Default()
	router.GET("/polls", handlers.GetAllPolls)
	// router.GET("/polls/:id", handlers.GetAPollByID)
	router.POST("/polls", handlers.PostAPoll)
	// router.POST("/polls/:pollID/:optionID", handlers.VoteOnAPoll)
	// router.DELETE("/polls/:id", handlers.DeleteAPollByID)
	// router.PATCH("/polls/:id", handlers.ModifyAPollByID)

	if err := router.Run("localhost:8080"); err != nil {
		log.Printf("Eroor running sercer: %v", err)
	}
}

func Connect() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/mongo?authSource=admin"))
	if err != nil {
		log.Fatal(err)
	}

	if err = Client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
