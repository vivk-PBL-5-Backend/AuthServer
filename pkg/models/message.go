package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	AuthorID      string             `json:"author_id" bson:"author_id"`
	DestinationID string             `json:"destination_id" bson:"destination_id"`
	Content       string             `json:"content" bson:"content"`
	Date          time.Time          `bson:"date"`
}
