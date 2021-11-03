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
	chat := new(models.Chat)

	_, err := r.messageDB.InsertOne(ctx, message)
	if err != nil {
		log.Errorf("error on inserting message: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	if err = r.chatDB.FindOne(ctx, bson.M{"_id": message.DestinationID}).Decode(chat); err != nil {
		log.Errorf("error occured while getting chat from db: %s", err.Error())
		if err == mongo.ErrNoDocuments {
			return auth.ErrUserDoesNotExist
		}

		return err
	}

	chat.Messages = append(chat.Messages, message.ID)
	if _, err = r.chatDB.UpdateOne(ctx, bson.D{{"_id", message.DestinationID}}, bson.M{"$set": chat}); err != nil {
		log.Errorf("error on inserting message in chat: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}
