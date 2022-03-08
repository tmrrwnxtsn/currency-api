package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type Rate struct {
	ID             int       `json:"id"`
	FirstCurrency  string    `json:"first_currency"`
	SecondCurrency string    `json:"second_currency"`
	Value          float32   `json:"value"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

func (r *Rate) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.FirstCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&r.SecondCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&r.Value, validation.Required, validation.Min(0.0)),
		validation.Field(&r.LastUpdateTime, validation.Required, validation.Max(time.Now())),
	)
}
