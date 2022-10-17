package sqlstore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	queryString := `
	INSERT INTO users (email, encrypted_password)
	VALUES ($1, $2)
	RETURNING id;`
	return r.store.db.QueryRow(queryString, u.Email, u.EncryptedPassword).Scan(&u.ID)

}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	queryString := `
	SELECT id, email, encrypted_password
	FROM users
	WHERE email = $1;`
	u := &model.User{}
	if err := r.store.db.QueryRow(queryString, email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return u, nil
}
