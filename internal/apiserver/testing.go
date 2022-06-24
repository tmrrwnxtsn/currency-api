package apiserver

import (
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/config"
	"os"
	"testing"
)

func TestLogger(t *testing.T) *logrus.Logger {
	t.Helper()

	return logrus.New()
}

func TestConfig(t *testing.T) *config.Config {
	t.Helper()

	cfg := config.New()

	cfg.CurrencyAPIKey = os.Getenv("TEST_CURRENCY_API_KEY")
	if cfg.CurrencyAPIKey == "" {
		cfg.CurrencyAPIKey = "b2d66c60-9a47-11ec-bde0-db97a92aaea8"
	}

	cfg.DatabaseURL = os.Getenv("TEST_DATABASE_URL")
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "host=localhost dbname=currencyapi_test user=postgres password=qwerty sslmode=disable"
	}

	return cfg
}
