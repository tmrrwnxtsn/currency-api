package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"github.com/tmrrwnxtsn/currency-api/internal/store/sqlstore"
	"testing"
)

func TestRateRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("exchange_rate")

	st := sqlstore.New(db)

	testRate := model.TestRate(t)
	err := st.Rate().Create(testRate)

	assert.NotNil(t, testRate)
	assert.NoError(t, err)
}

func TestRateRepository_FindByFirstCurrency(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("exchange_rate")

	st := sqlstore.New(db)

	firstCurrency := "USD"
	_, err := st.Rate().FindByFirstCurrency(firstCurrency)
	assert.EqualError(t, err, store.ErrRowNotFound.Error())

	testRate := model.TestRate(t)
	testRate.FirstCurrency = firstCurrency
	_ = st.Rate().Create(testRate)

	rate, err := st.Rate().FindByFirstCurrency(firstCurrency)
	assert.NoError(t, err)
	assert.NotNil(t, rate)
}
