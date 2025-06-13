package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"polls/src/pckg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAUser(u User) error {
	u.ID = primitive.NewObjectID()
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		fmt.Printf("Eroor: %v\n", err)
		return err
	}
	u.Password = hashedPassword

	if !IsUserOkayToGoToDatabase(u) {
		return errors.New("invalid user to be logged in: same username or empty password")
	}

	collection := db.Client.Database("pollsdb").Collection("users")

	res, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		fmt.Printf("Eroor: %v\n", err)
		return err
	}
	id := res.InsertedID
	fmt.Printf("ID of inserted %s\n", id)
	return nil
}

func GetAUserByUsername(name string) (User, error) {
	users, _ := GetAllUsers()
	for _, u := range users {
		if u.Username == name {
			return u, nil
		}
	}
	return User{}, errors.New("user not found")
}

func GetAllUsers() ([]User, error) {

	collection := db.Client.Database("pollsdb").Collection("users")

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("[POLLS] this broke", err.Error())
	}

	defer cur.Close(context.Background())

	var users []User
	if err = cur.All(context.Background(), &users); err != nil {
		log.Println("error parsing polls", err)
		return nil, err
	}
	return users, nil
}
