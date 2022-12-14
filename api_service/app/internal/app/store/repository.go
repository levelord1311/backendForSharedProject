package store

import "github.com/levelord1311/backendForSharedProject/api_service/internal/app/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	CreateUserWithGoogle(user *model.User) error
	FindByEmail(string) (*model.User, error)
	FindByEmailGoogle(string) (*model.User, error)
	FindByUsername(string) (*model.User, error)
}

type EstateLotRepository interface {
	CreateEstateLot(lot *model.EstateLot) error
	GetFreshEstateLots() (*[]model.EstateLot, error)
}
