package sqlstore

import (
	"database/sql"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/app/model"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *model.User) error {
	if err := u.ValidateFields(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	queryString := `
	INSERT INTO users (username, email, encrypted_password)
	VALUES (?, ?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Username, u.Email, u.EncryptedPassword)
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	var timestamp *model.RawTime

	if err = r.store.db.QueryRow("SELECT created_at FROM users WHERE id = ?", retID).Scan(&timestamp); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	createdAt, err := timestamp.Time()
	if err != nil {
		return err
	}

	u.CreatedAt, u.RedactedAt = createdAt, createdAt
	u.ID = uint(retID)
	return nil

}

func (r *UserRepository) CreateUserWithGoogle(u *model.User) error {

	queryString := `
	INSERT INTO users (email, given_name, family_name)
	VALUES (?, ?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Email, u.GivenName, u.FamilyName)
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	var timestamp *model.RawTime

	if err = r.store.db.QueryRow("SELECT created_at FROM users WHERE id = ?", retID).Scan(&timestamp); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	createdAt, err := timestamp.Time()
	if err != nil {
		return err
	}

	u.CreatedAt, u.RedactedAt = createdAt, createdAt
	u.ID = uint(retID)
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

func (r *UserRepository) FindByEmailGoogle(email string) (*model.User, error) {
	queryString := `
	SELECT id, email
	FROM users
	WHERE email = ?;`
	u := &model.User{}
	if err := r.store.db.QueryRow(queryString, email).Scan(
		&u.ID,
		&u.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	queryString := `
	SELECT id, username, email, encrypted_password
	FROM users
	WHERE username = ?;`
	u := &model.User{}
	if err := r.store.db.QueryRow(queryString, username).Scan(
		&u.ID,
		&u.Username,
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
