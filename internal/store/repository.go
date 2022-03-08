package store

import "github.com/tmrrwnxtsn/currency-api/internal/model"

// RateRepository ...
type RateRepository interface {
	// Create ...
	Create(*model.Rate) error

	// Find ...
	Find(int) (*model.Rate, error)

	// FindByCurrencies ...
	FindByCurrencies(string, string) (*model.Rate, error)
}
