package repository

import (
	"errors"

	"github.com/tianhongw/tinyid/conn"
	"github.com/tianhongw/tinyid/pkg/log"
)

type Repository struct {
	DB     *conn.DB
	Logger log.Logger

	TinyId *TinyIdRepository
}

type Option func(*Repository)

func WithLogger(logger log.Logger) Option {
	return func(r *Repository) {
		r.Logger = logger
	}
}

func NewRepository(db *conn.DB, opts ...Option) (*Repository, error) {
	if db == nil {
		return nil, errors.New("nil database")
	}

	r := &Repository{
		DB:     db,
		Logger: log.DummyLogger,
	}

	for _, o := range opts {
		o(r)
	}

	r.TinyId = NewTinyIdRepository(r)

	return r, nil
}
