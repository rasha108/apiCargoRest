package sqlstore

import (
	"testing"

	"github.com/rasha108/apiCargoRest.git/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teadDown := TestDB(t, databaseURL)
	defer teadDown("users")

	s := New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Find(t *testing.T) {
	db, teadDown := TestDB(t, databaseURL)
	defer teadDown("users")

	s := New(db)
	u := model.TestUser(t)
	s.User().Create(u)
	u, err := s.User().Find(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teadDown := TestDB(t, databaseURL)
	defer teadDown("users")

	s := New(db)
	u := model.TestUser(t)
	s.User().Create(u)
	u, err := s.User().FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
