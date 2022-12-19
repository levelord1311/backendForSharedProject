package db

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"time"
)

var _ lot.Storage = &db{}

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

func (s *db) Create(ctx context.Context, lot *lot.Lot) (uint, error) {
	queryString := `
	INSERT INTO lots (
		user_id,
		type_of_estate,
		rooms,
		area,
		floor,
		max_floor,
		city,
		district,
		street,
		building,
		price
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		lot.CreatedByUserID,
		lot.TypeOfEstate,
		lot.Rooms,
		lot.Area,
		lot.Floor,
		lot.MaxFloor,
		lot.City,
		lot.District,
		lot.Street,
		lot.Building,
		lot.Price,
	)
	if err != nil {
		return 0, err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(retID), err
}

func (s *db) FindByLotID(ctx context.Context, id uint) (*lot.Lot, error) {
	l := &lot.Lot{}
	var createdAt, redactedAt *rawTime

	queryString := `
	SELECT *
	FROM lots 
	WHERE lot_id=?;`

	err := s.db.QueryRowContext(ctx, queryString, id).Scan(
		&l.ID,
		&l.CreatedByUserID,
		&l.TypeOfEstate,
		&l.Rooms,
		&l.Area,
		&l.Floor,
		&l.MaxFloor,
		&l.City,
		&l.District,
		&l.Street,
		&l.Building,
		&l.Price,
		&createdAt,
		&redactedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	l.CreatedAt, err = createdAt.time()
	if err != nil {
		return nil, err
	}

	l.RedactedAt, err = redactedAt.time()
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *db) FindByUserID(ctx context.Context, id uint) ([]*lot.Lot, error) {
	lotsByUser := make([]*lot.Lot, 0, 10)

	queryString := `
	SELECT *
	FROM lots
	WHERE user_id=?;`

	rows, err := s.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		l := &lot.Lot{}
		var createdAt, redactedAt *rawTime
		if err = rows.Scan(
			&l.ID,
			&l.CreatedByUserID,
			&l.TypeOfEstate,
			&l.Rooms,
			&l.Area,
			&l.Floor,
			&l.MaxFloor,
			&l.City,
			&l.District,
			&l.Street,
			&l.Building,
			&l.Price,
			&createdAt,
			&redactedAt,
		); err != nil {
			return nil, err
		}
		l.CreatedAt, err = createdAt.time()
		if err != nil {
			return nil, err
		}

		l.RedactedAt, err = redactedAt.time()
		if err != nil {
			return nil, err
		}
		lotsByUser = append(lotsByUser, l)
	}
	if err = rows.Err(); err != nil {
		return lotsByUser, err
	}
	return lotsByUser, nil
}

func (s *db) Update(ctx context.Context, lot *lot.Lot) error {
	queryString := `
	UPDATE lots
	SET price=?
	WHERE lot_id=? AND user_id=?;`
	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, lot.Price, lot.ID, lot.CreatedByUserID)
	if err != nil {
		return err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAff == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (s *db) Delete(ctx context.Context, lotID, userID uint) error {
	queryString := `
	DELETE
	FROM lots 
	WHERE lot_id=? AND user_id=?;`
	stmt, err := s.db.PrepareContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, lotID, userID)
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAff == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
