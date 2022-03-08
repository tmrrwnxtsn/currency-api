package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
	"net/http"
	"strings"
	"time"
)

const (
	currencyAPIURLTemplate        = "https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=%s"
	ctxKeyRequestID        ctxKey = iota
)

type ctxKey int8

type server struct {
	router         *mux.Router
	logger         *logrus.Logger
	store          store.Store
	currencyAPIKey string
}

type freeAPIResponse struct {
	Query query              `json:"query"`
	Data  map[string]float32 `json:"data"`
}

type query struct {
	BaseCurrency string `json:"base_currency"`
}

func newServer(store store.Store, currencyAPIKey string) *server {
	srv := &server{
		router:         mux.NewRouter(),
		logger:         logrus.New(),
		store:          store,
		currencyAPIKey: currencyAPIKey,
	}

	srv.configureRouter()

	srv.logger.Info("starting API server")

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
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
	))

	s.router.HandleFunc("/api/create", s.handleCreateRate()).Methods("POST")
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

		logger.Infof(
			"completed with %d %s in %v",
			rw.code, http.StatusText(rw.code),
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
			s.error(w, http.StatusConflict, fmt.Errorf("rate for '%s'-'%s' already exists", req.FirstCurrency, req.SecondCurrency))
			return
		}

		netClient := http.Client{
			Timeout: time.Second * 10,
		}

		url := fmt.Sprintf(currencyAPIURLTemplate, s.currencyAPIKey, req.FirstCurrency)

		res, err := netClient.Get(url)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			s.error(w, res.StatusCode, fmt.Errorf("response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body))
			return
		}

		freeAPIRes := &freeAPIResponse{}
		if err = json.NewDecoder(res.Body).Decode(freeAPIRes); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}

		exchangeRateValue, ok := freeAPIRes.Data[req.SecondCurrency]
		if !ok {
			s.error(w, http.StatusUnprocessableEntity, fmt.Errorf("currency '%s' not found", req.SecondCurrency))
			return
		}

		rate = &model.Rate{
			FirstCurrency:  freeAPIRes.Query.BaseCurrency,
			SecondCurrency: req.SecondCurrency,
			Value:          exchangeRateValue,
			LastUpdateTime: time.Now(),
		}

		if err = s.store.Rate().Create(rate); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
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
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
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
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
