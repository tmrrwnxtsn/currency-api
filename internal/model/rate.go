package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type Rate struct {
	ID             int       `json:"id" example:"1"`
	FirstCurrency  string    `json:"first_currency" example:"RUB"`
	SecondCurrency string    `json:"second_currency" example:"USD"`
	Value          float32   `json:"value" example:"75.4"`
	LastUpdateTime time.Time `json:"last_update_time" example:"2019-11-09T21:21:46+00:00"`
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
