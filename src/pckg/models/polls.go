package models

import (
	"context"
	"fmt"
	"log"
	"polls/src/pckg/db"

	//"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

// var collection *mongo.Collection

// func InitCollection(c *mongo.Collection) {
// 	collection = c
//}

func GetAllPolls() ([]Poll, error) {
	log.Println("beginning")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := db.Client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("[POLLS] yes", err.Error())
		defer cancel()
		return []Poll{}, err
	}
	log.Println("beginning2")

	collection := db.Client.Database("pollsdb").Collection("polls")

	defer cancel()

	cur, err := collection.Find(ctx, bson.M{})
	log.Println("beginning3")

	if err != nil {
		log.Println("[POLLS] this broke", err.Error())
	}

	defer cur.Close(ctx)
	log.Println("beginning4")

	var polls []Poll
	if err = cur.All(ctx, &polls); err != nil {
		log.Println("error parsing polls", err)
		return nil, err
	}
	log.Println("here:", polls)
	return polls, nil
}

// func EditAPollByID(newPoll Poll) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	filter := bson.D{{Key: "_id", Value: newPoll.ID}}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "question", Value: newPoll.Question},
// 			{Key: "options", Value: newPoll.AnswerOptions},
// 			{Key: "closed", Value: newPoll.IsClosed},
// 		}},
// 	}
// 	collection.UpdateOne(ctx, filter, update)
// }

func AddAPoll(p Poll) {
	p.ID = primitive.NewObjectID()
	collection := db.Client.Database("pollsdb").Collection("polls")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, p)
	if err != nil {
		fmt.Printf("Eroor: %v\n", err)
		return
	}
	id := res.InsertedID
	fmt.Printf("ID of inserted %s\n", id)
}

// func DeleteAPollByID(id primitive.ObjectID) bool {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	resp,err:=collection.DeleteOne(ctx, bson.M{"_id": id})
// 	if err!=ni

// 	for i := 0; i < len(polls); i++ {
// 		if polls[i].ID == id {
// 			polls = slices.Delete(polls, i, i+1)
// 			return true
// 		}
// 	}
// 	return false
// }
