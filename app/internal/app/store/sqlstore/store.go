package sqlstore

import (
	"backendForSharedProject/internal/app/store"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db                  *sql.DB
	userRepository      *UserRepository
	estateLotRepository *EstateLotRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{s}
	}
	return s.userRepository
}

func (s *Store) EstateLot() store.EstateLotRepository {
	if s.estateLotRepository == nil {
		s.estateLotRepository = &EstateLotRepository{s}
	}
	return s.estateLotRepository
}
