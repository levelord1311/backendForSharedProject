package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "testDB:testPW@tcp(127.0.0.1:3306)/test_db"
	}
	os.Exit(m.Run())
}
