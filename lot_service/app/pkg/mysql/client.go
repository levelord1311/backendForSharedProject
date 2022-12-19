package mysql

import (
	"database/sql"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
)

func NewClient(logger logging.Logger, databaseURL string) (*sql.DB, error) {
	logger.Println("Opening DB...")
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
