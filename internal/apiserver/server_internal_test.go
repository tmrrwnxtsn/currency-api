package apiserver

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store/teststore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleCreateRate(t *testing.T) {
	srv := newServer(TestConfig(t), teststore.New(), TestLogger(t))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"first_currency":  "USD",
				"second_currency": "RUB",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "already exists",
			payload: map[string]string{
				"first_currency":  "USD",
				"second_currency": "RUB",
			},
			expectedCode: http.StatusConflict,
		},
		{
			name: "identical currencies",
			payload: map[string]string{
				"first_currency":  "USD",
				"second_currency": "USD",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"first_currency":  "dollar",
				"second_currency": "ruble",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			_ = json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/rate", b)

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleConvertCurrency(t *testing.T) {
	srv := newServer(TestConfig(t), teststore.New(), TestLogger(t))

	r := model.TestRate(t)
	_ = srv.store.Rate().Create(r)

	testCases := []struct {
		name         string
		payload      map[string]string
		expectedCode int
	}{
		{
			name:         "missing params",
			payload:      map[string]string{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid value",
			payload: map[string]string{
				"currency_from": r.FirstCurrency,
				"currency_to":   r.SecondCurrency,
				"value":         "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid currency",
			payload: map[string]string{
				"currency_from": "dollar",
				"currency_to":   "ruble",
				"value":         "10",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "valid",
			payload: map[string]string{
				"currency_from": r.FirstCurrency,
				"currency_to":   r.SecondCurrency,
				"value":         "10",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/convert", nil)

			q := req.URL.Query()
			for pkey, pvalue := range tc.payload {
				q.Add(pkey, pvalue)
			}
			req.URL.RawQuery = q.Encode()

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
