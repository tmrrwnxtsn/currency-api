package apiserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/currency-api/internal/config"
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store/sqlstore"
	"net/http"
	"time"
)

func Start(cfg *config.Config) error {
	db, err := newDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)

	logger := logrus.New()

	netClient := http.Client{
		Timeout: time.Second * 10,
	}

	go func(store *sqlstore.Store, logger *logrus.Logger) {
		for range time.Tick(time.Minute * time.Duration(cfg.UpdateInterval)) {
			rates, err := store.Rate().FindAll()
			if err != nil {
				logger.Errorf("error occurred while updating rates: %s", err.Error())
				continue
			}

			for _, r := range rates {
				// request to the external currency conversion API for the rates
				rateApiRequestUrl := fmt.Sprintf(
					"https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=%s",
					cfg.CurrencyAPIKey,
					r.FirstCurrency,
				)

				res, err := netClient.Get(rateApiRequestUrl)
				if err != nil {
					logger.Errorf("error occurred while updating '%s'-'%s' rate: %s", r.FirstCurrency, r.SecondCurrency, err.Error())
					continue
				}

				if res.StatusCode > 299 {
					logger.Errorf("error occurred while updating '%s'-'%s' rate: response failed with status code: %d and body: %s", r.FirstCurrency, r.SecondCurrency, res.StatusCode, res.Body)
					continue
				}

				response := &currencyApiResponse{}
				if err = json.NewDecoder(res.Body).Decode(response); err != nil {
					logger.Errorf("error occurred while updating rates: %s", err.Error())
					continue
				}

				rateUpd := &model.Rate{
					ID:             r.ID,
					FirstCurrency:  r.FirstCurrency,
					SecondCurrency: r.SecondCurrency,
					Value:          response.Data[r.SecondCurrency],
					LastUpdateTime: time.Now(),
				}

				if err = store.Rate().Update(rateUpd); err != nil {
					logger.Errorf("error occurred while updating '%s'-'%s' rate: %s", r.FirstCurrency, r.SecondCurrency, fmt.Errorf("response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body).Error())
					continue
				}

				res.Body.Close()

				logger.Infof("'%s'-'%s' rate was successfully updated!", r.FirstCurrency, r.SecondCurrency)

				time.Sleep(5 * time.Second)
			}
		}
	}(store, logger)

	srv := newServer(cfg, store, logger)

	return http.ListenAndServe(cfg.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
