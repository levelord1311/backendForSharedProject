package model

import (
	"testing"
)

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Username: "username_example",
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestLot(t *testing.T) *EstateLot {
	t.Helper()

	return &EstateLot{
		TypeOfEstate: "квартира",
		Rooms:        3,
		Area:         51,
		Floor:        6,
		MaxFloor:     9,
		City:         "Урюпинск",
		District:     "Ленинский",
		Street:       "Пушкина",
		Building:     "Колотушкина",
		Price:        55000,
	}
}
