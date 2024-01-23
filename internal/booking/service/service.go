package service

import "time"

// service implements user.Service.
type service struct {
	pgStore PGStore
	timeNow func() time.Time
}

// New creates a new service.
func New(pgStore PGStore) (*service, error) {
	s := &service{
		pgStore: pgStore,
		timeNow: time.Now,
	}

	return s, nil
}
