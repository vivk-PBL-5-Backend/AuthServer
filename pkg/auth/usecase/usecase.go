package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/repository/mongo"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"time"
)

type Authorizer struct {
	userRepo    auth.Repository
	chatRepo    auth.Repository
	messageRepo *mongo.MessageRepository

	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthorizer(userRepo auth.Repository, chatRepo auth.Repository, messageRepo *mongo.MessageRepository,
	hashSalt string, signingKey []byte, expireDuration time.Duration) *Authorizer {
	return &Authorizer{
		userRepo:       userRepo,
		chatRepo:       chatRepo,
		messageRepo:    messageRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *Authorizer) SignUp(ctx context.Context, user *models.User) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	result := a.userRepo.Insert(ctx, user)
	if result != nil {
		return result
	}

	_ = a.chatRepo.Insert(ctx, user)
	return result
}

func (a *Authorizer) SignIn(ctx context.Context, user *models.User) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.Get(ctx, user.Username, user.Password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.Username,
	})

	return token.SignedString(a.signingKey)
}

func (a *Authorizer) Send(ctx context.Context, message *models.Message) error {
	result := a.messageRepo.Send(ctx, message)
	if result != nil {
		return result
	}

	return result
}
