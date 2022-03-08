package model

import (
	"testing"
	"time"
)

func TestRate(t *testing.T) *Rate {
	t.Helper()

	return &Rate{
		FirstCurrency:  "USD",
		SecondCurrency: "RUB",
		Value:          121.41,
		LastUpdateTime: time.Now(),
	}
}
