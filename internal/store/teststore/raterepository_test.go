package teststore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"github.com/tmrrwnxtsn/currency-api/internal/store/teststore"
	"testing"
	"time"
)

func TestRateRepository_Create(t *testing.T) {
	st := teststore.New()

	testCases := []struct {
		name    string
		r       func() *model.Rate
		isValid bool
	}{
		{
			name: "valid",
			r: func() *model.Rate {
				return model.TestRate(t)
			},
			isValid: true,
		},
		{
			name: "invalid",
			r: func() *model.Rate {
				return &model.Rate{
					FirstCurrency:  "dollar",
					SecondCurrency: "ruble",
					Value:          -1,
					LastUpdateTime: time.Now().Add(time.Second * 10),
				}
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := st.Rate().Create(tc.r())
			if tc.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
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

func TestRateRepository_FindAll(t *testing.T) {
	st := teststore.New()

	rates := []*model.Rate{
		{
			FirstCurrency:  "USD",
			SecondCurrency: "RUB",
			Value:          75.1,
			LastUpdateTime: time.Now(),
		},
		{
			FirstCurrency:  "EUR",
			SecondCurrency: "USD",
			Value:          1.1,
			LastUpdateTime: time.Now(),
		},
		{
			FirstCurrency:  "BRL",
			SecondCurrency: "CAD",
			Value:          31.51,
			LastUpdateTime: time.Now(),
		},
	}

	for _, rate := range rates {
		err := st.Rate().Create(rate)
		assert.NoError(t, err)
	}

	rates, err := st.Rate().FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(rates))
}

func TestRateRepository_Update(t *testing.T) {
	st := teststore.New()

	r := model.TestRate(t)
	err := st.Rate().Create(r)
	assert.NoError(t, err)

	rUpd := &model.Rate{
		ID:             r.ID,
		FirstCurrency:  "EUR",
		SecondCurrency: "USD",
		Value:          1.1,
		LastUpdateTime: r.LastUpdateTime,
	}

	err = st.Rate().Update(rUpd)
	assert.NoError(t, err)

	rFind, err := st.Rate().Find(r.ID)
	assert.NoError(t, err)
	assert.Equal(t, rUpd.FirstCurrency, rFind.FirstCurrency)
	assert.Equal(t, rUpd.SecondCurrency, rFind.SecondCurrency)
	assert.Equal(t, rUpd.Value, rFind.Value)
}
