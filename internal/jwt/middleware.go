package jwt

import (
	"github.com/golang-jwt/jwt"
	"net/http"
)

func VerifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					_, err := w.Write([]byte("You're Unauthorized!"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err2 := w.Write([]byte("You're unauthorized due to error parsing JWT"))
				if err2 != nil {
					return
				}
			}
			if token.Valid {
				endpointHandler(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("You are unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You are unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}

	})
}
