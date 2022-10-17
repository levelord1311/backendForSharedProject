package store

import "github.com/levelord1311/HTTP-REST-API.git/internal/app/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(string) (*model.User, error)
}
