package lot

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

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

type CreateLotDTO struct {
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
}

type UpdateLotDTO struct {
	ID              uint `json:"id"`
	CreatedByUserID uint `json:"created_by_user_id"`
	Price           int  `json:"price"`
}

func NewLot(dto *CreateLotDTO) *Lot {
	return &Lot{
		CreatedByUserID: dto.CreatedByUserID,
		TypeOfEstate:    dto.TypeOfEstate,
		Rooms:           dto.Rooms,
		Area:            dto.Area,
		Floor:           dto.Floor,
		MaxFloor:        dto.MaxFloor,
		City:            dto.City,
		District:        dto.District,
		Street:          dto.Street,
		Building:        dto.Building,
		Price:           dto.Price,
	}
}

func UpdatedLot(dto *UpdateLotDTO) *Lot {
	return &Lot{
		ID:              dto.ID,
		CreatedByUserID: dto.CreatedByUserID,
		Price:           dto.Price,
	}
}

func (l *Lot) ValidateFields() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.CreatedByUserID, validation.Required),
		validation.Field(&l.TypeOfEstate, validation.Required, validation.In(
			"квартира",
			"дом")),
		validation.Field(&l.Rooms, validation.Max(6)),
		validation.Field(&l.Area, validation.Required),
		validation.Field(&l.Floor, validation.Required, validation.Max(163)),
		validation.Field(&l.MaxFloor, validation.Required, validation.Max(163)),
		validation.Field(&l.City, validation.Required),
		validation.Field(&l.District, validation.Required),
		validation.Field(&l.Street, validation.Required),
		validation.Field(&l.Building, validation.Required),
		validation.Field(&l.Price, validation.Required),
	)
}

func (dto *UpdateLotDTO) ValidateFields() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.ID, validation.Required),
		validation.Field(&dto.CreatedByUserID, validation.Required),
		validation.Field(&dto.Price, validation.Required))
}
