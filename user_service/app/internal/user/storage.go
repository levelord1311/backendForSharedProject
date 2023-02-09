package user

import (
	"context"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/models"
)

type Storage interface {
	Create(ctx context.Context, user *models.User) (uint, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}
