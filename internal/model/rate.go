package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type Rate struct {
	ID             int
	FirstCurrency  string
	SecondCurrency string
	Value          int
	LastUpdateTime time.Time
}

func (r *Rate) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.FirstCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&r.SecondCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&r.Value, validation.Required, validation.Min(0)),
		validation.Field(&r.LastUpdateTime, validation.Required),
	)
}
