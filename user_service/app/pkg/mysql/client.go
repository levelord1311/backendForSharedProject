package mysql

import (
	"context"
	"database/sql"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/logging"
	"time"
)

func NewClient(ctx context.Context, logger logging.Logger, databaseURL string) (*sql.DB, error) {
	logger.Println("Opening DB...")
	// TODO это правильное использование контекста? он дальше никуда не пробрасывается
	// cancel() сработает по истечению 10 секунд?
	_, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

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