package teststore

import (
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
)

var _ store.RateRepository = (*RateRepository)(nil)

type RateRepository struct {
	store *Store
	rates map[string]*model.Rate
}

func (r *RateRepository) Create(rate *model.Rate) error {
	if err := rate.Validate(); err != nil {
		return err
	}

	r.rates[rate.FirstCurrency] = rate
	rate.ID = len(r.rates)

	return nil
}

func (r *RateRepository) FindByFirstCurrency(firstCurrency string) (*model.Rate, error) {
	rate, ok := r.rates[firstCurrency]
	if !ok {
		return nil, store.ErrRowNotFound
	}

	return rate, nil
}
