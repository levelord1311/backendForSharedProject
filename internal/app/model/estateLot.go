package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type EstateLot struct {
	ID           uint   `json:"id"`
	TypeOfEstate string `json:"type_of_estate"`
	Rooms        int    `json:"rooms"`
	Area         int    `json:"area"`
	Floor        int    `json:"floor"`
	MaxFloor     int    `json:"max_floor"`
	City         string `json:"city"`
	District     string `json:"district"`
	Street       string `json:"street"`
	Building     string `json:"building"`
	Price        int    `json:"price"`
	CreatedAt    time.Time
	RedactedAt   time.Time
}

func (l *EstateLot) ValidateLotFields() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.TypeOfEstate, validation.Required, validation.In("квартира", "дом")),
		validation.Field(&l.Rooms, validation.Required, validation.Max(6)),
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
