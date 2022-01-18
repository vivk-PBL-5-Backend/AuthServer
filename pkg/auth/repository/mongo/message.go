package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/aes"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/cipheradapter"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/filereader"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	rsa2 "github.com/vivk-PBL-5-Backend/AuthServer/pkg/rsa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type MessageRepository struct {
	messageDB *mongo.Collection
	chatDB    *mongo.Collection

	rsaCipher cipheradapter.ICipher
	aesCipher cipheradapter.ICipher
}

func NewMessageRepository(db *mongo.Database, chatCollection string, messageCollection string) *MessageRepository {
	messageKey := filereader.ReadFile(viper.GetString("aes.message-key"))
	ivKey := filereader.ReadFile(viper.GetString("aes.iv-key"))

	aesCipher := aes.New([]byte(messageKey), []byte(ivKey))

	privateKey := filereader.ReadFile(viper.GetString("rsa.message-key"))
	rsaCipher := rsa2.New([]byte(privateKey))

	return &MessageRepository{
		messageDB: db.Collection(messageCollection),
		chatDB:    db.Collection(chatCollection),

		rsaCipher: rsaCipher,
		aesCipher: aesCipher,
	}
}

func (r *MessageRepository) Send(ctx context.Context, message *models.Message) error {
	message.AuthorID = r.aesCipher.Encrypt(message.AuthorID)
	message.DestinationID = r.aesCipher.Encrypt(message.DestinationID)
	message.Content = r.rsaCipher.Encrypt(message.Content)

	_, err := r.messageDB.InsertOne(ctx, message)
	if err != nil {
		log.Errorf("error on inserting message: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}

func (r *MessageRepository) Get(ctx context.Context, user string, companion string) ([]models.Message, error) {
	userID := r.aesCipher.Encrypt(user)
	companionID := r.aesCipher.Encrypt(companion)

	messages := make([]models.Message, 0)

	cursor, err := r.messageDB.Find(ctx, bson.M{"author_id": userID, "destination_id": companionID})
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

	cursor, err = r.messageDB.Find(ctx, bson.M{"author_id": companionID, "destination_id": userID})
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

	return r.messagesDecrypt(messages), nil
}

func (r *MessageRepository) messagesDecrypt(messages []models.Message) []models.Message {
	for i, _ := range messages {
		messages[i].AuthorID = strings.TrimSpace(r.aesCipher.Decrypt(messages[i].AuthorID))
		messages[i].DestinationID = strings.TrimSpace(r.aesCipher.Decrypt(messages[i].DestinationID))
		messages[i].Content = strings.TrimSpace(r.rsaCipher.Decrypt(messages[i].Content))
	}
	return messages
}

func (r *MessageRepository) Remove(ctx context.Context, user string, companion string) error {

	userID := r.aesCipher.Encrypt(user)
	companionID := r.aesCipher.Encrypt(companion)

	cursor, err := r.messageDB.Find(ctx, bson.M{"destination_id": companionID, "author_id": userID})
	if err == nil {
		for cursor.Next(ctx) {
			var elem models.Message
			err = cursor.Decode(&elem)
			if err != nil {
				log.Errorf("error on get message: %s", err.Error())
				return err
			}

			_, err = r.messageDB.DeleteOne(ctx, bson.M{"_id": elem.ID})
			if err != nil {
				log.Errorf("error on get message: %s", err.Error())
				return err
			}
		}
	}

	return nil
}
