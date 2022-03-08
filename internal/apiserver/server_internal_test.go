package apiserver

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store/teststore"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var currencyAPIKey string

func TestMain(m *testing.M) {
	currencyAPIKey = os.Getenv("TEST_CURRENCY_API_KEY")
	if currencyAPIKey == "" {
		currencyAPIKey = "b2d66c60-9a47-11ec-bde0-db97a92aaea8"
	}

	os.Exit(m.Run())
}

func TestServer_HandleCreateRate(t *testing.T) {
	srv := newServer(teststore.New(), currencyAPIKey)

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
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/create", b)

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleConvertCurrency(t *testing.T) {
	srv := newServer(teststore.New(), currencyAPIKey)

	r := model.TestRate(t)
	_ = srv.store.Rate().Create(r)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"currency_from": r.FirstCurrency,
				"currency_to":   r.SecondCurrency,
				"value":         "10",
			},
			expectedCode: http.StatusOK,
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
				"value":           "-1",
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodGet, "/api/convert", b)

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
