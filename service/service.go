package service

import (
	"errors"

	"github.com/tianhongw/tinyid/config"
	"github.com/tianhongw/tinyid/pkg/log"
	"github.com/tianhongw/tinyid/repository"
)

type Service struct {
	Repo   *repository.Repository
	Conf   *config.Config
	Logger log.Logger
}

type Option func(*Service)

func WithLogger(logger log.Logger) Option {
	return func(s *Service) {
		s.Logger = logger
	}
}

func NewService(repo *repository.Repository, opts ...Option) (*Service, error) {
	if repo == nil {
		return nil, errors.New("nil repository")
	}

	srv := &Service{
		Repo:   repo,
		Logger: log.DummyLogger,
	}

	for _, o := range opts {
		o(srv)
	}

	return srv, nil
}
