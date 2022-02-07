package store

import "github.com/tmrrwnxtsn/currency-api/internal/model"

type RateRepository struct {
	store *Store
}

func (r *RateRepository) Create(rate *model.Rate) (*model.Rate, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO exchange_rate (first_currency, second_currency, rate_value, last_update_time) VALUES ($1, $2, $3, $4) RETURNING id",
		rate.FirstCurrency, rate.SecondCurrency, rate.RateValue, rate.LastUpdateTime,
	).Scan(&rate.ID); err != nil {
		return nil, err
	}

	return rate, nil
}

func (r *RateRepository) FindByFirstCurrency(firstCurrency string) (*model.Rate, error) {
	rate := &model.Rate{}
	if err := r.store.db.QueryRow(
		"SELECT id, first_currency, second_currency, rate_value, last_update_time FROM exchange_rate WHERE first_currency = $1",
		firstCurrency,
	).Scan(&rate.ID, &rate.FirstCurrency, &rate.SecondCurrency, &rate.RateValue, &rate.LastUpdateTime); err != nil {
		return nil, err
	}

	return rate, nil
}
