package sqlstore

import (
	"database/sql"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
)

var _ store.RateRepository = (*RateRepository)(nil)

type RateRepository struct {
	store *Store
}

func (r *RateRepository) Create(rate *model.Rate) error {
	if err := rate.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO exchange_rate (first_currency, second_currency, value, last_update_time) VALUES ($1, $2, $3, $4) RETURNING id",
		rate.FirstCurrency, rate.SecondCurrency, rate.Value, rate.LastUpdateTime,
	).Scan(&rate.ID)
}

func (r *RateRepository) FindByFirstCurrency(firstCurrency string) (*model.Rate, error) {
	rate := &model.Rate{}
	if err := r.store.db.QueryRow(
		"SELECT id, first_currency, second_currency, value, last_update_time FROM exchange_rate WHERE first_currency = $1",
		firstCurrency,
	).Scan(&rate.ID, &rate.FirstCurrency, &rate.SecondCurrency, &rate.Value, &rate.LastUpdateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRowNotFound
		}

		return nil, err
	}

	return rate, nil
}
