package sqlstore

import (
	"github.com/jmoiron/sqlx"

	"github.com/rasha108/apiCargoRest.git/internal/app/db"
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

func (s *Store) User() db.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
