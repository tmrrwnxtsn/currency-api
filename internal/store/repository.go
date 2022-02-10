package store

import "github.com/tmrrwnxtsn/currency-api/internal/model"

// RateRepository ...
type RateRepository interface {
	// Create ...
	Create(rate *model.Rate) error

	// FindByFirstCurrency ...
	FindByFirstCurrency(firstCurrency string) (*model.Rate, error)
}
