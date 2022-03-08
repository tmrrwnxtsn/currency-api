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

func TestRateRepository_Find(t *testing.T) {
	st := teststore.New()
	r1 := model.TestRate(t)
	_ = st.Rate().Create(r1)

	r2, err := st.Rate().Find(r1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, r2)
}

func TestRateRepository_FindByCurrencies(t *testing.T) {
	st := teststore.New()
	r1 := model.TestRate(t)
	_, err := st.Rate().FindByCurrencies(r1.FirstCurrency, r1.SecondCurrency)
	assert.EqualError(t, err, store.ErrRowNotFound.Error())

	_ = st.Rate().Create(r1)
	r2, err := st.Rate().FindByCurrencies(r1.FirstCurrency, r1.SecondCurrency)
	assert.NoError(t, err)
	assert.NotNil(t, r2)
}
