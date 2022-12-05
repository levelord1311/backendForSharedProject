package apiserver

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store/teststore"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO replace hardcoded jwtKey as env var
func TestServer_HandleDefaultPage(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s := newServer(teststore.New(), []byte("jwtKeyExample"))
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TODO replace hardcoded jwtKey as env var
func TestServer_HandleUsersCreate(t *testing.T) {
	s := newServer(teststore.New(), []byte("jwtKeyExample"))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"username": "username_example",
				"email":    "user@example.org",
				"password": "passwordExample",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, "/users/create", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// TODO replace hardcoded jwtKey as env var
func TestServer_HandleJWTCreate(t *testing.T) {
	u := model.TestUser(t)
	store := teststore.New()
	store.User().CreateUser(u)
	s := newServer(store, []byte("jwtKeyExample"))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid username",
			payload: map[string]string{
				"login":    u.Username,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid email",
			payload: map[string]string{
				"login":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid login",
			payload: map[string]string{
				"login":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"login":    u.Username,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "empty password",
			payload: map[string]string{
				"login":    u.Username,
				"password": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty login",
			payload: map[string]string{
				"login":    "",
				"password": u.Password,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, "/auth", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// TODO replace hardcoded jwtKey as env var
func TestServer_HandleEstateLotsCreate(t *testing.T) {
	s := newServer(teststore.New(), []byte("jwtKeyExample"))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid type_of_estate 1",
			payload: map[string]interface{}{
				"type_of_estate": "квартира",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "valid type_of_estate 2",
			payload: map[string]interface{}{
				"type_of_estate": "дом",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "all fields empty",
			payload: map[string]interface{}{
				"type_of_estate": "",
				"rooms":          0,
				"area":           0,
				"floor":          0,
				"max_floor":      0,
				"city":           "",
				"district":       "",
				"street":         "",
				"building":       "",
				"price":          0,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid type_of_estate field",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid payload",
			payload:      "not a map",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "rooms field empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          0,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "rooms field less than 0",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          -1,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "rooms field larger than 6",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          7,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "area field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           0,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "area field less than 0",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           -1,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "floor field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          0,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "floor field is less than 0",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          -1,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "floor field is larger than 163",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          164,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "max floor field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      0,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "max floor field is less than 0",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      -1,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "max floor field is larger than 163",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      164,
				"city":           "Магадан",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "city field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "",
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "city field not a string",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           17,
				"district":       "Северный",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "district field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "",
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "district field not a string",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       12,
				"street":         "Владимирская",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "street field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         "",
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "street field not a string",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         42,
				"building":       "8",
				"price":          12000,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "building field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         "Пушкина",
				"building":       "",
				"price":          12000,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "building field not a string",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         "",
				"building":       8,
				"price":          12000,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "price field is empty",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         "",
				"building":       "8",
				"price":          0,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "price field not an int",
			payload: map[string]interface{}{
				"type_of_estate": "invalid",
				"rooms":          2,
				"area":           51,
				"floor":          6,
				"max_floor":      9,
				"city":           "Магадан",
				"district":       "Ленинский",
				"street":         "",
				"building":       "8",
				"price":          "12000",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest(http.MethodPost, "/estate_lots/create", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
