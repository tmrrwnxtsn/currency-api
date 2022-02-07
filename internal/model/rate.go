package model

import "time"

type Rate struct {
	ID             int
	FirstCurrency  string
	SecondCurrency string
	RateValue      int
	LastUpdateTime time.Time
}
