package teststore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	userRepository      *UserRepository
	estateLotRepository *EstateLotRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[string]*model.User),
		}
	}

	return s.userRepository
}

func (s *Store) EstateLot() store.EstateLotRepository {
	if s.estateLotRepository == nil {
		s.estateLotRepository = &EstateLotRepository{
			store: s,
			lots:  make(map[uint]*model.EstateLot)}
	}
	return s.estateLotRepository
}
