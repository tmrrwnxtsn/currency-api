package teststore

import (
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/model"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/store"
)

var _ store.RateRepository = (*RateRepository)(nil)

type RateRepository struct {
	store *Store
	rates map[int]*model.Rate
}

func (r *RateRepository) Create(rate *model.Rate) error {
	if err := rate.Validate(); err != nil {
		return err
	}

	rate.ID = len(r.rates) + 1
	r.rates[rate.ID] = rate

	return nil
}

func (r *RateRepository) Find(id int) (*model.Rate, error) {
	rate, ok := r.rates[id]
	if !ok {
		return nil, store.ErrRowNotFound
	}

	return rate, nil
}

func (r *RateRepository) FindByCurrencies(firstCurrency, secondCurrency string) (*model.Rate, error) {
	for _, rate := range r.rates {
		if rate.FirstCurrency == firstCurrency && rate.SecondCurrency == secondCurrency {
			return rate, nil
		}
	}

	return nil, store.ErrRowNotFound
}

func (r *RateRepository) FindAll() ([]*model.Rate, error) {
	var rates []*model.Rate

	for id, rate := range r.rates {
		rate.ID = id
		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *RateRepository) Update(rate *model.Rate) error {
	_, err := r.Find(rate.ID)
	if err != nil {
		return err
	}

	r.rates[rate.ID] = rate

	return nil
}
