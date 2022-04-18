package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tmrrwnxtsn/currency-api/docs"
	"github.com/tmrrwnxtsn/currency-api/internal/config"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const ctxKeyRequestID ctxKey = iota

var (
	errMissingRequiredParams = errors.New("one or more required parameters are missing")
	errWrongValueParam       = errors.New("parameter 'value' is wrong")
	errIdenticalCurrencies   = errors.New("the exchange rate should contain information about different currencies")
)

type ctxKey int8

type server struct {
	config *config.Config
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(config *config.Config, store store.Store, logger *logrus.Logger) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
		config: config,
	}

	srv.configureRouter()

	srv.logger.Info("API server started")

	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedMethods([]string{"*"}),
	))

	s.router.HandleFunc("/api/v1/rate", s.handleCreateRate()).Methods("POST")
	s.router.HandleFunc("/api/v1/convert", s.handleConvertCurrency()).Methods("GET")

	// swagger documentation
	s.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
}

func (s server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

type createRateQuery struct {
	FirstCurrency  string `json:"first_currency" example:"RUB"`
	SecondCurrency string `json:"second_currency" example:"USD"`
}

// handleCreateRate godoc
// @Summary      Create an exchange rate
// @Description  create a record of the exchange rate between two currencies
// @Tags         rate
// @Accept       json
// @Produce      json
// @Param        input  body      createRateQuery  true  "An exchange rate information"
// @Success      201    {object}  model.Rate       "Ok"
// @Failure      400    {object}  errorResponse    "Missing parameters or invalid payload"
// @Failure      409    {object}  errorResponse    "An exchange rate record with these currencies already exists"
// @Failure      422    {object}  errorResponse    "Invalid parameters"
// @Failure      500    {object}  errorResponse
// @Router       /rate [post]
func (s *server) handleCreateRate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &createRateQuery{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		if req.FirstCurrency == "" || req.SecondCurrency == "" {
			s.error(w, http.StatusBadRequest, errMissingRequiredParams)
			return
		}

		if req.FirstCurrency == req.SecondCurrency {
			s.error(w, http.StatusUnprocessableEntity, errIdenticalCurrencies)
			return
		}

		rate, _ := s.store.Rate().FindByCurrencies(strings.ToUpper(req.FirstCurrency), strings.ToUpper(req.SecondCurrency))
		if rate != nil {
			s.error(w, http.StatusConflict, fmt.Errorf("the exchange rate record for %s-%s already exists", req.FirstCurrency, req.SecondCurrency))
			return
		}

		res, err := getExchangeRates(s.config.CurrencyAPIKey, req.FirstCurrency)
		if err != nil {
			s.error(w, http.StatusUnprocessableEntity, fmt.Errorf("error occurred while getting exchange rates for the currency %s: %s", req.FirstCurrency, err.Error()))
			return
		}

		exchangeRateValue, ok := res.Data[req.SecondCurrency]
		if !ok {
			s.error(w, http.StatusUnprocessableEntity, fmt.Errorf("currency %s not found", req.SecondCurrency))
			return
		}

		rate = &model.Rate{
			FirstCurrency:  res.Query.BaseCurrency,
			SecondCurrency: req.SecondCurrency,
			Value:          exchangeRateValue,
			LastUpdateTime: time.Now(),
		}

		if err = s.store.Rate().Create(rate); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusCreated, rate)
	}
}

type convertCurrencyQuery struct {
	CurrencyFrom string  `json:"currency_from" example:"RUB"`
	CurrencyTo   string  `json:"currency_to" example:"RUB"`
	Value        float32 `json:"value,string" example:"123.321"`
}

type convertCurrencyResponse struct {
	Query            convertCurrencyQuery `json:"query"`
	ConversionResult float32              `json:"conversion_result,string" example:"123.321"`
	LastUpdateTime   time.Time            `json:"last_update_time" example:"2019-11-09T21:21:46+00:00"`
}

// handleConvertCurrency godoc
// @Summary      Currency conversion
// @Description  convert the value from one currency to another according to the exchange rate
// @Tags         other
// @Accept       json
// @Produce      json
// @Param        currency_from  query     string                   true  "The currency whose value will be converted to another currency"
// @Param        currency_to    query     string                   true  "The currency to which the value from the first currency will be converted"
// @Param        value          query     number                   true  "The value that will be converted from one currency to another"
// @Success      200            {object}  convertCurrencyResponse  "Ok"
// @Failure      400            {object}  errorResponse            "Missing parameters"
// @Failure      404            {object}  errorResponse            "There is no record of the exchange rate"
// @Failure      422            {object}  errorResponse            "Invalid parameters"
// @Failure      500            {object}  errorResponse
// @Router       /convert [get]
func (s *server) handleConvertCurrency() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if !(q.Has("currency_from") && q.Has("currency_to") && q.Has("value")) {
			s.error(w, http.StatusBadRequest, errMissingRequiredParams)
			return
		}

		valueFloat64, err := strconv.ParseFloat(q.Get("value"), 32)
		if err != nil {
			s.error(w, http.StatusUnprocessableEntity, errWrongValueParam)
			return
		}

		req := &convertCurrencyQuery{
			CurrencyFrom: q.Get("currency_from"),
			CurrencyTo:   q.Get("currency_to"),
			Value:        float32(valueFloat64),
		}

		rate, err := s.store.Rate().FindByCurrencies(req.CurrencyFrom, req.CurrencyTo)
		if err != nil {
			if err == store.ErrRowNotFound {
				s.error(w, http.StatusNotFound, err)
				return
			}

			s.error(w, http.StatusInternalServerError, err)
			return
		}

		res := &convertCurrencyResponse{
			Query:            *req,
			ConversionResult: req.Value * rate.Value,
			LastUpdateTime:   rate.LastUpdateTime,
		}

		s.respond(w, http.StatusOK, res)
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func (s *server) error(w http.ResponseWriter, statusCode int, err error) {
	s.respond(w, statusCode, errorResponse{err.Error()})
}

func (s *server) respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
