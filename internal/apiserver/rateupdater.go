package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/config"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/model"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/store/sqlstore"
	"net/http"
	"time"
)

type rateUpdater struct {
	config *config.Config
	store  *sqlstore.Store
	logger *logrus.Logger
}

func newRateUpdater(config *config.Config, store *sqlstore.Store, logger *logrus.Logger) *rateUpdater {
	return &rateUpdater{
		config: config,
		store:  store,
		logger: logger,
	}
}

// Start ...
func (u *rateUpdater) Start() {
	for range time.Tick(time.Minute * time.Duration(u.config.UpdateInterval)) {
		rates, err := u.store.Rate().FindAll()
		if err != nil {
			u.logger.Errorf("error occurred while getting rates from the db: %s", err.Error())
			continue
		}

		for _, rate := range rates {
			response, err := getExchangeRates(u.config.CurrencyAPIKey, rate.FirstCurrency)
			if err != nil {
				u.logger.Errorf("error occurred while getting the rate info for the currency %s: %s", rate.FirstCurrency, err.Error())
				continue
			}

			rateUpd := &model.Rate{
				ID:             rate.ID,
				FirstCurrency:  rate.FirstCurrency,
				SecondCurrency: rate.SecondCurrency,
				Value:          response.Data[rate.SecondCurrency],
				LastUpdateTime: time.Now(),
			}

			if err = u.store.Rate().Update(rateUpd); err != nil {
				u.logger.Errorf("error occurred while updating %s-%s rate: %s", rate.FirstCurrency, rate.SecondCurrency, err.Error())
				continue
			}

			u.logger.Infof("%s-%s rate was successfully updated!", rate.FirstCurrency, rate.SecondCurrency)

			time.Sleep(5 * time.Second)
		}
	}
}

type getExchangeRatesQuery struct {
	BaseCurrency string `json:"base_currency"`
}

type getExchangeRatesResponse struct {
	Query getExchangeRatesQuery `json:"query"`
	Data  map[string]float32    `json:"data"`
}

func getExchangeRates(currencyApiKey, baseCurrency string) (*getExchangeRatesResponse, error) {
	// request to the external currency conversion API for the rates
	rateApiRequestUrl := fmt.Sprintf(
		"https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=%s",
		currencyApiKey,
		baseCurrency,
	)

	res, err := http.Get(rateApiRequestUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf(
			"error occurred while updating rate for the currency %s: response failed with status code: %d and\nbody: %s\n",
			baseCurrency, res.StatusCode, res.Body,
		)
	}

	response := &getExchangeRatesResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
