package teststore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"github.com/tmrrwnxtsn/currency-api/internal/store/teststore"
	"testing"
)

func TestRateRepository_Create(t *testing.T) {
	st := teststore.New()

	testRate := model.TestRate(t)
	err := st.Rate().Create(testRate)

	assert.NotNil(t, testRate)
	assert.NoError(t, err)
}

func TestRateRepository_FindByFirstCurrency(t *testing.T) {
	st := teststore.New()

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
