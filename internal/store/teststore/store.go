package teststore

import (
	"github.com/tmrrwnxtsn/currency-api/internal/model"
	"github.com/tmrrwnxtsn/currency-api/internal/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	rateRepository *RateRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Rate() store.RateRepository {
	if s.rateRepository != nil {
		return s.rateRepository
	}

	s.rateRepository = &RateRepository{
		store: s,
		rates: make(map[string]*model.Rate),
	}

	return s.rateRepository
}
