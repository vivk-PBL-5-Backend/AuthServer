package mongo

import (
	"context"
	"crypto/rsa"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/aes"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
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

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	cipher aes.ICipher
}

func NewMessageRepository(db *mongo.Database, chatCollection string, messageCollection string) *MessageRepository {
	messageKey := filereader.ReadFile(viper.GetString("aes.message-key"))
	ivKey := filereader.ReadFile(viper.GetString("aes.iv-key"))

	cipher := aes.New([]byte(messageKey), []byte(ivKey))

	privateKeyPath := viper.GetString("rsa.message-key")
	publicKey, privateKey := rsa2.GenerateKeyPair(privateKeyPath)

	return &MessageRepository{
		messageDB: db.Collection(messageCollection),
		chatDB:    db.Collection(chatCollection),

		publicKey:  publicKey,
		privateKey: privateKey,

		cipher: cipher,
	}
}

func (r *MessageRepository) Send(ctx context.Context, message *models.Message) error {
	message.AuthorID = r.cipher.Encrypt(message.AuthorID)
	message.DestinationID = r.cipher.Encrypt(message.DestinationID)
	message.Content = rsa2.Encrypt(r.publicKey, message.Content)

	_, err := r.messageDB.InsertOne(ctx, message)
	if err != nil {
		log.Errorf("error on inserting message: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}

func (r *MessageRepository) Get(ctx context.Context, user string, companion string) ([]models.Message, error) {
	userID := r.cipher.Encrypt(user)
	companionID := r.cipher.Encrypt(companion)

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
		messages[i].AuthorID = strings.TrimSpace(r.cipher.Decrypt(messages[i].AuthorID))
		messages[i].DestinationID = strings.TrimSpace(r.cipher.Decrypt(messages[i].DestinationID))
		messages[i].Content = strings.TrimSpace(rsa2.Decrypt(r.privateKey, messages[i].Content))
	}
	return messages
}
