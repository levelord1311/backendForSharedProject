package sqlstore

import (
	"database/sql"
	"backendForSharedProject/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{s}
	}
	return s.userRepository
}
