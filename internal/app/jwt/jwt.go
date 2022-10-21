package jwt

import (
	"backendForSharedProject/internal/app/model"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateJWT(u *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["sub"] = u.ID

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("You're Unauthorized due to invalid method"))
				}
				return "", nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("You're unauthorized due to error parsing JWT"))
				return
			}
			if token.Valid {
				endpointHandler(w, r)
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("You are unauthorized due to invalid token"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are unauthorized due to No token in the header"))
			return
		}
	})
}
