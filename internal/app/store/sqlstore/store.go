package sqlstore

import (
	"github.com/jmoiron/sqlx"

	"github.com/rasha108/apiCargoRest.git/internal/app/store"
)

type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
