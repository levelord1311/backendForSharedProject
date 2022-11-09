package sqlstore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"database/sql"
	"time"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	queryString := `
	INSERT INTO users (username, email, encrypted_password, created_at)
	VALUES (?, ?, ?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Username, u.Email, u.EncryptedPassword, time.Now())
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint(retID)
	return nil

}

func (r *UserRepository) CreateUserWithGoogle(u *model.User) error {

	queryString := `
	INSERT INTO users (email, given_name, family_name, created_at)
	VALUES (?, ?, ?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Email, u.GivenName, u.FamilyName, time.Now())
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

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

func (r *UserRepository) CreateEstateLot(lot *model.EstateLot) error {
	queryString := `
	INSERT INTO estate_lots (
	                  type_of_estate,
	                  rooms,
	                  area,
	                  floor,
	                  max_floor,
	                  city,
	                  district,
	                  street,
	                  building,
	                  price,
	                  created_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	stmt, err := r.store.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
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
		time.Now(),
	)
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	lot.ID = uint(retID)
	return nil

}

func (r *UserRepository) GetAllEstateLots() (*[]model.EstateLot, error) {

	//нужно будет ограничить количество выводимых лотов и соответственно изменить размер создаваемого в памяти слайса.
	lots := make([]model.EstateLot, 50)

	rows, err := r.store.db.Query("SELECT * FROM estate_lots")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		i := 0
		err := rows.Scan(&lots[i])
		if err != nil {
			return nil, err
		}
		i++
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &lots, nil
}
