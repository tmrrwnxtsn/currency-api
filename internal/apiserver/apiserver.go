package apiserver

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/config"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/store/sqlstore"
	"net/http"
)

// Start ...
func Start(cfg *config.Config) error {
	db, err := newDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	logger := logrus.New()

	updater := newRateUpdater(cfg, store, logger)
	go updater.Start()

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
