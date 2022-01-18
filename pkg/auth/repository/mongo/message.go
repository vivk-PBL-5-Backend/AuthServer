package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	messageDB *mongo.Collection
	chatDB    *mongo.Collection
}

func NewMessageRepository(db *mongo.Database, chatCollection string, messageCollection string) *MessageRepository {
	return &MessageRepository{
		messageDB: db.Collection(messageCollection),
		chatDB:    db.Collection(chatCollection),
	}
}

func (r *MessageRepository) Send(ctx context.Context, message *models.Message) error {
	_, err := r.messageDB.InsertOne(ctx, message)
	if err != nil {
		log.Errorf("error on inserting message: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}

func (r *MessageRepository) Get(ctx context.Context, destination string, author string) ([]models.Message, error) {
	messages := make([]models.Message, 0)

	cursor, err := r.messageDB.Find(ctx, bson.M{"destination_id": author, "author_id": destination})
	if err == nil {
		for cursor.Next(ctx) {
			var elem models.Message
			err = cursor.Decode(&elem)
			if err != nil {
				log.Errorf("error on get message: %s", err.Error())
				return nil, auth.ErrUserAlreadyExists
			}

			messages = append(messages, elem)
		}
	}

	cursor, err = r.messageDB.Find(ctx, bson.M{"destination_id": destination, "author_id": author})
	if err == nil {
		for cursor.Next(ctx) {
			var elem models.Message
			err = cursor.Decode(&elem)
			if err != nil {
				log.Errorf("error on get message: %s", err.Error())
				return nil, auth.ErrUserAlreadyExists
			}

			_, err = r.messageDB.DeleteOne(ctx, bson.M{"_id": elem.ID})
			if err != nil {
				log.Errorf("error on get message: %s", err.Error())
				return nil, auth.ErrUserAlreadyExists
			}

			messages = append(messages, elem)
		}
	}

	return messages, nil
}
