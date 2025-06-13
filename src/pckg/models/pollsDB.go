package models

import (
	"context"
	"fmt"
	"log"
	"polls/src/pckg/db"

	//"slices"
	// "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/v2/mongo/readpref"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

// var collection *mongo.Collection

// func InitCollection(c *mongo.Collection) {
// 	collection = c
//}

func GetAllPolls() ([]Poll, error) {
	// log.Println("beginning")
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// if err := db.Client.Ping(ctx, readpref.Primary()); err != nil {
	// 	log.Println("[POLLS] yes", err.Error())
	// 	defer cancel()
	// 	return []Poll{}, err
	// }

	collection := db.Client.Database("pollsdb").Collection("polls")

	// defer cancel()

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("[POLLS] this broke", err.Error())
	}

	defer cur.Close(context.Background())
	log.Println("beginning4")

	var polls []Poll
	if err = cur.All(context.Background(), &polls); err != nil {
		log.Println("error parsing polls", err)
		return nil, err
	}
	// log.Println("here:", polls)
	return polls, nil
}

func (p Poll) EditAPollByID(newPoll Poll) error {
	//newPoll.ID = primitive.NewObjectID()

	collection := db.Client.Database("pollsdb").Collection("polls")

	filter := bson.M{"_id": p.ID}
	update := bson.M{
		"$set": bson.M{
			"question": newPoll.Question,
			"options":  newPoll.AnswerOptions,
			"closed":   newPoll.IsClosed,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func AddAPoll(p Poll) {
	collection := db.Client.Database("pollsdb").Collection("polls")

	res, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		fmt.Printf("Eroor: %v\n", err)
		return
	}
	id := res.InsertedID
	fmt.Printf("ID of inserted %s\n", id)
}

func DeleteAPollByID(id primitive.ObjectID) bool {
	collection := db.Client.Database("pollsdb").Collection("polls")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Print("Error deleteing a poll: ", err)
		return false
	}
	return true
}
