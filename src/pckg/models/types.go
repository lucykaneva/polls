package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Poll struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedUserID primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	Question      string             `bson:"question" json:"question"`
	AnswerOptions []Option           `bson:"options" json:"options"`
	IsClosed      bool               `bson:"closed" json:"closed"`
}

type Option struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content string             `bson:"content" json:"content"`
	Votes   []Vote             `bson:"votes" json:"votes"`
}

type Vote struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"userId" json:"userId"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}
