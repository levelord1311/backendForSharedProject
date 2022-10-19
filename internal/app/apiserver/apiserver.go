package apiserver

import (
	"backendForSharedProject/internal/app/store/sqlstore"
	"database/sql"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func StartHTTP(config *Config) error {
	log.Println("Starting HTTP server...")
	return http.ListenAndServe(config.BindAddr, http.HandlerFunc(redirectToTls))

}

func StartMainHTTP(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	log.Println("Starting main HTTP server...")
	s := newServer(store, sessionStore)
	return http.ListenAndServe(config.BindAddr, s)

}

func StartTLS(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	log.Println("Starting TLS server...")
	s := newServer(store, sessionStore)
	return http.ListenAndServeTLS(config.TLSAddr, config.Cert, config.Key, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	log.Println("Opening DB...")
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	log.Println("Establishing DB connection...")
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}
