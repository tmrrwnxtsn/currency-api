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
		Value:          85,
		LastUpdateTime: time.Now(),
	}
}
