package sqlstore_test

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=currencyapi_test user=postgres password=qwerty sslmode=disable"
	}

	os.Exit(m.Run())
}
