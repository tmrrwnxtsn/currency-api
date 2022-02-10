package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"testing"
)

func TestRate_Validate(t *testing.T) {
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
			name: "empty first currency",
			r: func() *model.Rate {
				testRate := model.TestRate(t)
				testRate.FirstCurrency = ""
				return testRate
			},
			isValid: false,
		},
		{
			name: "invalid first currency format",
			r: func() *model.Rate {
				testRate := model.TestRate(t)
				testRate.FirstCurrency = "dollar"
				return testRate
			},
			isValid: false,
		},
		{
			name: "empty second currency",
			r: func() *model.Rate {
				testRate := model.TestRate(t)
				testRate.SecondCurrency = ""
				return testRate
			},
			isValid: false,
		},
		{
			name: "invalid second currency format",
			r: func() *model.Rate {
				testRate := model.TestRate(t)
				testRate.SecondCurrency = "ruble"
				return testRate
			},
			isValid: false,
		},
		{
			name: "negative rate value",
			r: func() *model.Rate {
				testRate := model.TestRate(t)
				testRate.Value = -1
				return testRate
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.r().Validate())
			} else {
				assert.Error(t, tc.r().Validate())
			}
		})
	}
}
