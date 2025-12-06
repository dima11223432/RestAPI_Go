package teststore_test

import (
	"RestApi/internal/app/model"
	"RestApi/internal/app/store"
	teststore "RestApi/internal/app/store/testStore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.NewStore()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.NewStore()
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
