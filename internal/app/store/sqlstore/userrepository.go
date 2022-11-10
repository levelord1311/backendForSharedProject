package sqlstore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *model.User) error {
	if err := u.ValidateUserFields(); err != nil {
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
	if err := lot.ValidateLotFields(); err != nil {
		return err
	}

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
	                  price
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

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
	)
	if err != nil {
		return err
	}
	retID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	var timestamp *model.RawTime

	if err = r.store.db.QueryRow("SELECT created_at FROM estate_lots WHERE id = ?", retID).Scan(&timestamp); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	createdAt, err := timestamp.Time()
	if err != nil {
		return err
	}

	lot.CreatedAt, lot.RedactedAt = createdAt, createdAt
	lot.ID = uint(retID)
	return nil

}

func (r *UserRepository) GetAllEstateLots() (*[]model.EstateLot, error) {

	//нужно будет ограничить количество выводимых лотов и соответственно изменить размер создаваемого в памяти слайса.
	lots := make([]model.EstateLot, 20)

	queryString := `
	SELECT (
		    id,
    		type_of_estate,
	        rooms,
		    area,
	    	floor,
	        max_floor,
	        district,
	        street,
	        building,
	        price,
	       	redacted_at
	) FROM estate_lots
	ORDER BY redacted_at DESC;`

	rows, err := r.store.db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		i := 0
		var redactedTS *model.RawTime
		err := rows.Scan(
			&lots[i].ID,
			&lots[i].TypeOfEstate,
			&lots[i].Rooms,
			&lots[i].Area,
			&lots[i].Floor,
			&lots[i].MaxFloor,
			&lots[i].District,
			&lots[i].Street,
			&lots[i].Building,
			&lots[i].Price,
			&redactedTS,
		)

		if err != nil {
			return nil, err
		}

		redactedAt, err := redactedTS.Time()
		if err != nil {
			return nil, err
		}

		lots[i].RedactedAt = redactedAt
		i++
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &lots, nil
}
