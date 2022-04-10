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

var testCurrencyAPIUrl string

func TestMain(m *testing.M) {
	testCurrencyAPIKey := os.Getenv("TEST_CURRENCY_API_KEY")
	if testCurrencyAPIKey == "" {
		testCurrencyAPIKey = "b2d66c60-9a47-11ec-bde0-db97a92aaea8"
	}

	testCurrencyAPIUrl = "https://freecurrencyapi.net/api/v2/latest?apikey=" + testCurrencyAPIKey + "&base_currency=%s"

	os.Exit(m.Run())
}

func TestServer_HandleCreateRate(t *testing.T) {
	srv := newServer(teststore.New(), testCurrencyAPIUrl)

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

			req, _ := http.NewRequest(http.MethodPost, "/api/create", b)

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleConvertCurrency(t *testing.T) {
	srv := newServer(teststore.New(), testCurrencyAPIUrl)

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
			req, _ := http.NewRequest(http.MethodGet, "/api/convert", nil)

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
