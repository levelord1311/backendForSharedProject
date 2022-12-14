package db

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/user"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/logging"
	"time"
)

var _ user.Storage = &db{}

type db struct {
	db     *sql.DB
	logger logging.Logger
}

func NewStorage(storage *sql.DB, logger logging.Logger) *db {
	return &db{
		db:     storage,
		logger: logger,
	}
}

type rawTime []byte

func (t *rawTime) time() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(*t))
}

func (s *db) Create(ctx context.Context, u *user.User) (uint, error) {

	queryString := `
	INSERT INTO users (username, email, encrypted_password)
	VALUES (?, ?, ?);`

	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, u.Username, u.Email, u.EncryptedPassword)
	if err != nil {
		return 0, err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(retID), nil
}

func (s *db) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	var createdAt, redactedAt *rawTime

	queryString := `
	SELECT id, username, email, 
	IFNULL(given_name, ""),
	IFNULL(family_name, ""),
	created_at, redacted_at
	FROM users 
	WHERE email=?;`

	err := s.db.QueryRowContext(ctx, queryString, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.GivenName,
		&u.FamilyName,
		&createdAt,
		&redactedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	u.CreatedAt, err = createdAt.time()
	if err != nil {
		return nil, err
	}

	u.RedactedAt, err = redactedAt.time()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *db) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	u := &user.User{}
	var createdAt, redactedAt *rawTime

	queryString := `
	SELECT id, username, email, 
	IFNULL(given_name, ""),
	IFNULL(family_name, ""),
	created_at, redacted_at
	FROM users 
	WHERE username=?;`

	err := s.db.QueryRowContext(ctx, queryString, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.GivenName,
		&u.FamilyName,
		&createdAt,
		&redactedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	u.CreatedAt, err = createdAt.time()
	if err != nil {
		return nil, err
	}

	u.RedactedAt, err = redactedAt.time()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *db) FindByID(ctx context.Context, id uint) (*user.User, error) {

	u := &user.User{}
	var createdAt, redactedAt *rawTime

	queryString := `
	SELECT id, username, email, 
	IFNULL(given_name, ""),
	IFNULL(family_name, ""),
	created_at, redacted_at
	FROM users 
	WHERE id=?;`

	err := s.db.QueryRowContext(ctx, queryString, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.GivenName,
		&u.FamilyName,
		&createdAt,
		&redactedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	u.CreatedAt, err = createdAt.time()
	if err != nil {
		return nil, err
	}

	u.RedactedAt, err = redactedAt.time()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *db) Update(ctx context.Context, user *user.User) error {
	queryString := `
	UPDATE users
	SET encrypted_password=?
	WHERE id=?;`
	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.EncryptedPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *db) Delete(ctx context.Context, id uint) error {
	queryString := `
	DELETE
	FROM users 
	WHERE id=?;`
	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
