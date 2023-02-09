package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	expectedModel     = &models.User{}
	exampleUserReturn = &models.User{
		ID:                1234,
		Username:          "testUser",
		Email:             "test@email.org",
		Password:          "",
		EncryptedPassword: "",
		GivenName:         "Name",
		FamilyName:        "Surname",
		CreatedAt:         time.Time{},
		RedactedAt:        time.Time{},
	}
)

type stubService struct {
	err error
}

func (s *stubService) GetByID(ctx context.Context, id int) (*models.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	return exampleUserReturn, nil
}

func TestHandler_GetUser(t *testing.T) {
	t.Parallel()

	h := NewHandler(&stubService{})

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, singleUserURL, h.GetUser)

	cases := []struct {
		name           string
		idParam        string
		wantStatusCode int
		err            error
	}{
		{
			name:           "passing correct parameter",
			idParam:        "1",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "passing another correct parameter",
			idParam:        "1234",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "passing empty parameter",
			idParam:        "",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "passing parameter inconvertible to int",
			idParam:        "1a2b3c",
			wantStatusCode: http.StatusBadRequest,
			err:            apperror.ErrCantConvertID,
		},
		{
			name:           "user with such id is not found",
			idParam:        "666",
			wantStatusCode: http.StatusNotFound,
			err:            apperror.ErrNotFound,
		},
		{
			name:           "unexpected internal error",
			idParam:        "1",
			wantStatusCode: http.StatusTeapot,
			err:            apperror.ErrUnpredictedInternal,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

			h.service = &stubService{
				err: test.err,
			}

			url := fmt.Sprintf("%s/%s", usersURL, test.idParam)
			w := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, url, nil)

			router.ServeHTTP(w, request)

			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			defer response.Body.Close()
			assert.Equal(t, test.wantStatusCode, response.StatusCode)

			if test.err != nil {
				assert.Equal(t, test.err.Error(), string(body))
				return
			}

			receivedUser := &models.User{}
			switch test.wantStatusCode {
			case http.StatusNotFound:
				assert.Equal(t, expectedModel, receivedUser)
			default:
				err := json.Unmarshal(body, receivedUser)
				assert.NoError(t, err)
				assert.Equal(t, exampleUserReturn, receivedUser)
			}

		})
	}
}
