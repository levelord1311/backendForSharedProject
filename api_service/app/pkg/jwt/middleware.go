package jwt

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/config"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"net/http"
	"time"
)

func Middleware(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var claims *UserClaims
		logger := logging.GetLogger()
		key := []byte(config.GetConfig().JWT.Secret)

		logger.Debug("searching for \"Token\" in header...")
		if r.Header["Token"] == nil {
			err := errors.New("\"Token\" field has not been found in header")
			unauthorized(w, err)
			return
		}

		logger.Debug("\"Token\" field has been found, parsing...")

		token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				err := errors.New("wrong signing method, expected HMAC")
				unauthorized(w, err)
				return nil, err
			}
			logger.Debug("token signing method is correct")
			return key, nil
		})

		if err != nil {
			unauthorized(w, err)
			return
		}

		logger.Debug("checking if token is valid...")
		if !token.Valid {
			err = errors.New("token is not valid")
			unauthorized(w, err)
			return
		}

		logger.Debug("checking if token is expired...")
		if !claims.VerifyExpiresAt(time.Now(), true) {
			err = errors.New("token is expired")
			unauthorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_uuid", claims.ID)
		endpointHandler(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, err error) {
	logging.GetLogger().Error(err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
