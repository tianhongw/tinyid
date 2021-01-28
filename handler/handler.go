package handler

import (
	"errors"

	"github.com/tianhongw/tinyid/handler/option"
	v1 "github.com/tianhongw/tinyid/handler/v1"
	"github.com/tianhongw/tinyid/pkg/log"
	"github.com/tianhongw/tinyid/service"
)

type Handler struct {
	V1 *v1.HandlerV1
}

func NewHandler(services *service.Service, opt ...option.Option) (*Handler, error) {
	if services == nil {
		return nil, errors.New("nil service")
	}

	opts := &option.Options{
		Logger: log.DummyLogger,
	}

	for _, o := range opt {
		o(opts)
	}

	return &Handler{
		V1: v1.NewHandlerV1(services, opts),
	}, nil
}
