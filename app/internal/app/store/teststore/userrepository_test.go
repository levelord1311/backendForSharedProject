package teststore_test

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"backendForSharedProject/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().CreateUser(u))
}

func TestUserRepository_FindByEmail(t *testing.T) {

	s := teststore.New()
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().CreateUser(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByUsername(t *testing.T) {

	s := teststore.New()
	username := "username_example"
	_, err := s.User().FindByUsername(username)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Username = username
	if err = s.User().CreateUser(u); err != nil {
		t.Fatal(err)
	}

	u, err = s.User().FindByUsername(username)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
