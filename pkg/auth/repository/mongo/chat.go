package mongo

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
	db *mongo.Collection
}

func NewChatRepository(db *mongo.Database, collection string) *ChatRepository {
	return &ChatRepository{
		db: db.Collection(collection),
	}
}

func (r *ChatRepository) AddCompanion(ctx context.Context, userID string, companionID string) error {
	chat, err := r.findOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	companionIndex := -1
	for i, elem := range chat.Companions {
		if elem == companionID {
			companionIndex = i
			break
		}
	}

	if companionIndex != -1 {
		return errors.New("user[\"" + companionID + "\"] already in companion list!")
	}

	chat.Companions = append(chat.Companions, companionID)
	if _, err = r.db.UpdateOne(ctx, bson.D{{"_id", chat.Username}}, bson.M{"$set": chat}); err != nil {
		log.Errorf("error on inserting companion in chat: %s", err.Error())
		return err
	}

	return nil
}

func (r *ChatRepository) RemoveCompanion(ctx context.Context, userID string, companionID string) error {
	chat, err := r.findOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	companionIndex := -1
	for i, elem := range chat.Companions {
		if elem == companionID {
			companionIndex = i
			break
		}
	}

	if companionIndex == -1 {
		return errors.New("user[\"" + companionID + "\"] is not in companion list!")
	}

	chat.Companions = append(chat.Companions[:companionIndex], chat.Companions[companionIndex+1:]...)

	if _, err = r.db.UpdateOne(ctx, bson.D{{"_id", chat.Username}}, bson.M{"$set": chat}); err != nil {
		log.Errorf("error on inserting companion in chat: %s", err.Error())
		return err
	}

	return nil
}

func (r *ChatRepository) GetCompanions(ctx context.Context, userID string) ([]string, error) {
	chat, err := r.findOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return chat.Companions, nil
}

func (r *ChatRepository) findOrCreate(ctx context.Context, userID string) (*models.Chat, error) {
	chat := new(models.Chat)

	if err := r.db.FindOne(ctx, bson.M{"_id": userID}).Decode(chat); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Errorf("error occured while getting chat from db: %s", err.Error())
			return nil, err
		}

		chat.Username = userID
		chat.Companions = make([]string, 0)

		_, err = r.db.InsertOne(ctx, chat)
		if err != nil {
			log.Errorf("error on inserting chat: %s", err.Error())
			return nil, err
		}
	}

	return chat, nil
}
