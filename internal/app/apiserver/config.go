package apiserver

import (
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	TLSAddr     string `toml:"tls_addr"`
	Cert        string `toml:"cert"`
	Key         string `toml:"key"`
	DatabaseURL string `toml:"database_url"`
	JwtKey      []byte `toml:"jwt_key"`
}

var ErrMustBeSet = errors.New("must be set")

func NewConfig() (*Config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("$PORT %s", ErrMustBeSet)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("&DATABASE_URL %s", ErrMustBeSet)
	}

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if jwtKey == nil {
		return nil, fmt.Errorf("&JWT_KEY %s", ErrMustBeSet)
	}

	return &Config{
		BindAddr:    ":" + port,
		DatabaseURL: databaseURL,
		JwtKey:      jwtKey,
	}, nil
}

func NewGoogleConfig() (*oauth2.Config, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID %s", ErrMustBeSet)
	}

	googleClientSec := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSec == "" {
		return nil, fmt.Errorf("&GOOGLE_CLIENT_SECRET %s", ErrMustBeSet)
	}

	return &oauth2.Config{
		RedirectURL:  "https://backend-server-36962.herokuapp.com/auth/google/callback",
		ClientID:     googleClientID,
		ClientSecret: googleClientSec,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}, nil

}
