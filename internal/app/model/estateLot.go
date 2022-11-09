package model

import "time"

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
