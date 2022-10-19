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
	VALUES (?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Email, u.EncryptedPassword)
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(retID)
	return nil

}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	queryString := `
	SELECT id, email, encrypted_password
	FROM users
	WHERE email = ?;`
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
