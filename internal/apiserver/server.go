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
	s.router.HandleFunc("/api/rate", s.handleCreateRate()).Methods("POST")
	s.router.HandleFunc("/api/convert", s.handleConvertCurrency()).Methods("GET")
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

func (s *server) handleCreateRate() http.HandlerFunc {
	type request struct {
		FirstCurrency  string `json:"first_currency"`
		SecondCurrency string `json:"second_currency"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		rate, _ := s.store.Rate().FindByCurrencies(strings.ToUpper(req.FirstCurrency), strings.ToUpper(req.SecondCurrency))
		if rate != nil {
			s.error(w, http.StatusConflict, fmt.Errorf("rate for %s-%s already exists", req.FirstCurrency, req.SecondCurrency))
			return
		}

		res, err := getCurrencyRates(s.config.CurrencyAPIKey, req.FirstCurrency)
		if err != nil {
			s.error(w, http.StatusUnprocessableEntity, fmt.Errorf("error occurred while getting the rate info for the currency %s: %s", req.FirstCurrency, req.SecondCurrency))
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

func (s *server) handleConvertCurrency() http.HandlerFunc {
	type request struct {
		CurrencyFrom string  `json:"currency_from"`
		CurrencyTo   string  `json:"currency_to"`
		Value        float32 `json:"value,string"`
	}

	type response struct {
		Query            request   `json:"query"`
		ConvertingResult float32   `json:"converting_result,string"`
		LastUpdateTime   time.Time `json:"last_update_time"`
	}

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

		req := &request{
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

		res := &response{
			Query:            *req,
			ConvertingResult: req.Value * rate.Value,
			LastUpdateTime:   rate.LastUpdateTime,
		}

		s.respond(w, http.StatusOK, res)
	}
}

func (s *server) error(w http.ResponseWriter, statusCode int, err error) {
	s.respond(w, statusCode, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
