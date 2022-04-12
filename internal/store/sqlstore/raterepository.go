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
		"INSERT INTO rate (first_currency, second_currency, value, last_update_time) VALUES ($1, $2, $3, $4) RETURNING id",
		rate.FirstCurrency, rate.SecondCurrency, rate.Value, rate.LastUpdateTime,
	).Scan(&rate.ID)
}

func (r *RateRepository) Find(id int) (*model.Rate, error) {
	rate := &model.Rate{}
	if err := r.store.db.QueryRow(
		"SELECT id, first_currency, second_currency, value, last_update_time FROM rate WHERE id = $1",
		id,
	).Scan(&rate.ID, &rate.FirstCurrency, &rate.SecondCurrency, &rate.Value, &rate.LastUpdateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRowNotFound
		}

		return nil, err
	}

	return rate, nil
}

func (r *RateRepository) FindByCurrencies(firstCurrency, secondCurrency string) (*model.Rate, error) {
	rate := &model.Rate{}
	if err := r.store.db.QueryRow(
		"SELECT id, first_currency, second_currency, value, last_update_time FROM rate WHERE first_currency = $1 AND second_currency = $2",
		firstCurrency, secondCurrency,
	).Scan(&rate.ID, &rate.FirstCurrency, &rate.SecondCurrency, &rate.Value, &rate.LastUpdateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRowNotFound
		}

		return nil, err
	}

	return rate, nil
}

func (r *RateRepository) FindAll() ([]*model.Rate, error) {
	var rates []*model.Rate

	rows, err := r.store.db.Query("SELECT id, first_currency, second_currency, value, last_update_time FROM rate")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rate := &model.Rate{}
		if err = rows.Scan(&rate.ID, &rate.FirstCurrency, &rate.SecondCurrency, &rate.Value, &rate.LastUpdateTime); err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rates, nil
}

func (r *RateRepository) Update(rate *model.Rate) error {
	if err := rate.Validate(); err != nil {
		return err
	}

	findRate, err := r.Find(rate.ID)
	if err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"UPDATE rate SET first_currency = $2, second_currency = $3, value = $4, last_update_time = $5 WHERE id = $1 RETURNING id",
		findRate.ID, rate.FirstCurrency, rate.SecondCurrency, rate.Value, rate.LastUpdateTime,
	).Scan(&rate.ID)
}
