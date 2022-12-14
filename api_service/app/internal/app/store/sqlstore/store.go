package sqlstore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/app/store"
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
