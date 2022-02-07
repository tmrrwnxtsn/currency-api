package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"testing"
	"time"
)

func TestRateRepository_Create(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("exchange_rate")

	rate, err := st.ExchangeRate().Create(&model.Rate{
		FirstCurrency:  "USD",
		SecondCurrency: "RUB",
		RateValue:      80,
		LastUpdateTime: time.Now(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, rate)
}

func TestRateRepository_FindByFirstCurrency(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("exchange_rate")

	firstCurrency := "USD"
	_, err := st.ExchangeRate().FindByFirstCurrency(firstCurrency)
	assert.Error(t, err)

	_, _ = st.ExchangeRate().Create(&model.Rate{
		FirstCurrency:  firstCurrency,
		SecondCurrency: "RUB",
		RateValue:      80,
		LastUpdateTime: time.Now(),
	})

	rate, err := st.ExchangeRate().FindByFirstCurrency(firstCurrency)
	assert.NoError(t, err)
	assert.NotNil(t, rate)
}
