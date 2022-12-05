package teststore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) CreateUserWithGoogle(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindByEmailGoogle(s string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {

	u, ok := r.users[username]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) CreateUser(u *model.User) error {
	if err := u.ValidateFields(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	r.users[u.Email] = u
	//TODO think something better or is it really acceptable
	r.users[u.Username] = u
	u.ID = uint(len(r.users))

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
