package auth

import (
	"context"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
)

type Repository interface {
	Insert(ctx context.Context, user *models.User) error
	Get(ctx context.Context, username, password string) (*models.User, error)
	Exist(ctx context.Context, userID string) error
}
