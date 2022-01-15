package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *ChatRepository) Insert(ctx context.Context, user *models.User) error {
	chat := new(models.Chat)
	chat.Username = user.Username
	chat.Messages = make([]primitive.ObjectID, 0)

	_, err := r.db.InsertOne(ctx, chat)
	if err != nil {
		log.Errorf("error on inserting user: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}

func (r *ChatRepository) Get(ctx context.Context, username, password string) (*models.User, error) {
	return nil, auth.ErrUserDoesNotExist
}
