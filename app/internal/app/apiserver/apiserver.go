package apiserver

import (
	"backendForSharedProject/internal/app/store/sqlstore"
	"backendForSharedProject/internal/config"
	"backendForSharedProject/pkg/logging"
	"database/sql"
	"net/http"
)

func StartMainHTTP(config *config.Config) error {
	logger := logging.GetLogger()
	db, err := newDB(config.Database.Url)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)

	logger.Info("Starting main HTTP server...")
	s := newServer(store, []byte(config.JWT.Secret))
	return http.ListenAndServe(config.Listen.Port, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	logger := logging.GetLogger()
	logger.Info("Opening DB...")
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	logger.Info("Establishing DB connection...")
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}
