package apiserver

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"log"
	"net/http"
	"time"
)

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

var GoogleOauthConfig, errGoogleConfigLoad = NewGoogleConfig()

func init() {
	if errGoogleConfigLoad != nil {
		log.Fatal("error loading Google Config:", errGoogleConfigLoad)
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(5 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

//func getGooglePublicKey(keyID string) (string, error) {
//	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
//	if err != nil {
//		return "", err
//	}
//
//	dat, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//
//	myResp := map[string]string{}
//	err = json.Unmarshal(dat, &myResp)
//	if err != nil {
//		return "", err
//	}
//
//	key, ok := myResp[keyID]
//	if !ok {
//		return "", errors.New("key not found")
//	}
//
//	return key, nil
//}
//
//func getClaimsFromGoogle(code string) (*GoogleClaims, error) {
//	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
//	if err != nil {
//		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
//	}
//
//	claims, err := verifyGoogleJWT(token.AccessToken, GoogleOauthConfig.ClientID)
//	if err != nil {
//		return nil, fmt.Errorf("claims extraction error: %s", err.Error())
//	}
//	return claims, nil
//}

//func verifyGoogleJWT(tokenString, clientID string) (*GoogleClaims, error) {
//	claimsStruct := GoogleClaims{}
//
//	token, err := jwt.ParseWithClaims(
//		tokenString,
//		&claimsStruct,
//		func(token *jwt.Token) (interface{}, error) {
//			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
//			if err != nil {
//				return "", err
//			}
//			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
//			if err != nil {
//				return "", err
//			}
//			return key, nil
//		})
//	if err != nil {
//		return nil, err
//	}
//	claims, ok := token.Claims.(*GoogleClaims)
//	if !ok {
//		return nil, errors.New("invalid Google JWT")
//	}
//
//	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
//		return nil, errors.New("issuer is invalid")
//	}
//
//	if claims.Audience != clientID {
//		return nil, errors.New("aud is invalid")
//	}
//
//	if claims.ExpiresAt < time.Now().UTC().Unix() {
//		return nil, errors.New("JWT is expired")
//	}
//
//	return claims, nil
//}
