package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	Username string               `json:"username" bson:"_id"`
	Messages []primitive.ObjectID `json:"messages" bson:"messages"`
}
