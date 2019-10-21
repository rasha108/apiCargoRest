package db

import (
	"github.com/google/uuid"
	"github.com/rasha108/apiCargoRest.git/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Organizations(uuid.UUID) (*model.Organizations, error)
}
