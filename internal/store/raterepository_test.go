package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"testing"
)

func TestRateRepository_Create(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("exchange_rate")

	testRate := model.TestRate(t)
	rate, err := st.ExchangeRate().Create(testRate)

	assert.NoError(t, err)
	assert.NotNil(t, rate)
}

func TestRateRepository_FindByFirstCurrency(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("exchange_rate")

	firstCurrency := "USD"
	_, err := st.ExchangeRate().FindByFirstCurrency(firstCurrency)
	assert.Error(t, err)

	testRate := model.TestRate(t)
	testRate.FirstCurrency = firstCurrency
	_, _ = st.ExchangeRate().Create(testRate)

	rate, err := st.ExchangeRate().FindByFirstCurrency(firstCurrency)
	assert.NoError(t, err)
	assert.NotNil(t, rate)
}
