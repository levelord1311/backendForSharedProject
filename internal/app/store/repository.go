package store

import "backendForSharedProject/internal/app/model"

type UserRepository interface {
	Create(user *model.User) error
	CreateWithGoogle(user *model.User) error
	FindByEmail(string) (*model.User, error)
}
