package sqlstore

import (
	"database/sql"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/store"

	_ "github.com/lib/pq"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db             *sql.DB
	rateRepository *RateRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Rate() store.RateRepository {
	if s.rateRepository != nil {
		return s.rateRepository
	}

	s.rateRepository = &RateRepository{
		store: s,
	}

	return s.rateRepository
}
