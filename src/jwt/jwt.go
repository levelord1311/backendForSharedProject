package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

// необходимо спрятать, т.к. любой у кого есть этот ключ может проводить авторизацию
var sampleSecretKey = []byte("SecretYouShouldHide")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaims(_ http.ResponseWriter, r *http.Request) (string, error) {
	if r.Header["Token"] != nil {
		tokenString := r.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's error with the signing method")
			}
			return sampleSecretKey, nil
		})

		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			return username, nil
		}
	}
	return "unable to extract claims", nil
}
