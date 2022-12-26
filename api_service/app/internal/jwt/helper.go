package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/user_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/config"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"strconv"
	"time"
)

// TODO implement refresh token

var _ Helper = &helper{}

type Helper interface {
	GenerateAccessToken(u *user_service.User) ([]byte, error)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email    string
	Username string
}

type helper struct {
	logger logging.Logger
}

func NewHelper(logger logging.Logger) Helper {
	return &helper{
		logger: logger,
	}
}

func (h *helper) GenerateAccessToken(u *user_service.User) ([]byte, error) {
	// TODO get user struct param from user_service
	key := []byte(config.GetConfig().JWT.Secret)

	userIDStr := strconv.Itoa(int(u.ID))

	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        userIDStr,
		},
		Email:    u.Email,
		Username: u.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return []byte(tokenString), nil
}
