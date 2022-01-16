package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/aes"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/filereader"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserRepository struct {
	db *mongo.Collection

	cipher aes.ICipher
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	authKey := filereader.ReadFile(viper.GetString("aes.auth-key"))
	ivKey := filereader.ReadFile(viper.GetString("aes.iv-key"))

	cipher := aes.New([]byte(authKey), []byte(ivKey))

	return &UserRepository{
		db: db.Collection(collection),

		cipher: cipher,
	}
}

func (r *UserRepository) Insert(ctx context.Context, user *models.User) error {
	user.Username = r.cipher.Encrypt(user.Username)

	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("error on inserting user: %s", err.Error())
		return auth.ErrUserAlreadyExists
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context, username, password string) (*models.User, error) {
	user := new(models.User)

	user.Username = r.cipher.Encrypt(username)
	if err := r.db.FindOne(ctx, bson.M{"_id": user.Username, "password": password}).Decode(user); err != nil {
		log.Errorf("error occured while getting user from db: %s", err.Error())
		if err == mongo.ErrNoDocuments {
			return nil, auth.ErrUserDoesNotExist
		}

		return nil, err
	}

	user.Username = strings.TrimSpace(r.cipher.Decrypt(user.Username))
	return user, nil
}
