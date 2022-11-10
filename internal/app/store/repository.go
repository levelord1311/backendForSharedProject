package store

import "backendForSharedProject/internal/app/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	CreateUserWithGoogle(user *model.User) error
	FindByEmail(string) (*model.User, error)
	FindByEmailGoogle(string) (*model.User, error)
	FindByUsername(string) (*model.User, error)
	CreateEstateLot(lot *model.EstateLot) error
	GetFreshEstateLots() (*[]model.EstateLot, error)
}
