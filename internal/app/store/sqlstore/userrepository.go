package sqlstore

import (
	"database/sql"
	"errors"

	"github.com/rasha108/apiCargoRest.git/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO public.create_users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM public.create_users where email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
	}

	return u, nil
}
