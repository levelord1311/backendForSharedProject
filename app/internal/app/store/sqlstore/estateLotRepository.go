package sqlstore

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"database/sql"
)

type EstateLotRepository struct {
	store *Store
}

func (r *EstateLotRepository) CreateEstateLot(lot *model.EstateLot) error {
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

func (r *EstateLotRepository) GetFreshEstateLots() (*[]model.EstateLot, error) {

	sliceSize := 100
	lots := make([]model.EstateLot, sliceSize)

	queryString := `
	SELECT *
	FROM estate_lots
	ORDER BY redacted_at DESC
	LIMIT 100;`

	rows, err := r.store.db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if i >= sliceSize-1 {
			lots = append(lots, make([]model.EstateLot, i*2)...)
		}
		var createdTS, redactedTS *model.RawTime
		err := rows.Scan(
			&lots[i].ID,
			&lots[i].TypeOfEstate,
			&lots[i].Rooms,
			&lots[i].Area,
			&lots[i].Floor,
			&lots[i].MaxFloor,
			&lots[i].City,
			&lots[i].District,
			&lots[i].Street,
			&lots[i].Building,
			&lots[i].Price,
			&createdTS,
			&redactedTS,
		)

		if err != nil {
			return nil, err
		}

		createdAt, err := createdTS.Time()
		if err != nil {
			return nil, err
		}

		redactedAt, err := redactedTS.Time()
		if err != nil {
			return nil, err
		}

		lots[i].CreatedAt, lots[i].RedactedAt = createdAt, redactedAt
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//truncate empty lots
	for k, lot := range lots {
		if (lot == model.EstateLot{}) {
			func() {
				lots = lots[:k]
			}()
			break
		}
	}

	return &lots, nil
}
