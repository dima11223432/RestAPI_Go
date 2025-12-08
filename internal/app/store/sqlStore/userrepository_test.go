package sqlstore_test

import (
	"RestApi/internal/app/model"
	"RestApi/internal/app/store"
	"RestApi/internal/app/store/sqlStore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	email := "user@example.org"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	assert.NoError(t, s.User().Create(u))

	found, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	u := model.TestUser(t)
	s.User().Create(u)
	UserToFind, err := s.User().Find(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, UserToFind)
}
