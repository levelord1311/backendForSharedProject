package lot_service

import "time"

type Lot struct {
	ID              uint   `json:"id"`
	CreatedByUserID uint   `json:"created_by_user_id"`
	TypeOfEstate    string `json:"type_of_estate"`
	Rooms           int    `json:"rooms"`
	Area            int    `json:"area"`
	Floor           int    `json:"floor"`
	MaxFloor        int    `json:"max_floor"`
	City            string `json:"city"`
	District        string `json:"district"`
	Street          string `json:"street"`
	Building        string `json:"building"`
	Price           int    `json:"price"`
	CreatedAt       time.Time
	RedactedAt      time.Time
}

// CreateLotDTO model info
// @Description lot information for registering in db.
type CreateLotDTO struct {
	CreatedByUserID uint   `json:"created_by_user_id"` // leave empty, value is taken from JWT
	TypeOfEstate    string `json:"type_of_estate"`     // required. either "квартира" or "дом"
	Rooms           int    `json:"rooms"`              // required. max - 6; 0 rooms means studio flat
	Area            int    `json:"area"`               // required.
	Floor           int    `json:"floor"`              // required. max - 163
	MaxFloor        int    `json:"max_floor"`          // required. max - 163
	City            string `json:"city"`               // required.
	District        string `json:"district"`           // required.
	Street          string `json:"street"`             // required.
	Building        string `json:"building"`           // required.
	Price           int    `json:"price"`              // required.
}

type UpdateLotDTO struct {
	ID              uint `json:"id"`
	CreatedByUserID uint `json:"created_by_user_id"`
	Price           int  `json:"price"`
}
