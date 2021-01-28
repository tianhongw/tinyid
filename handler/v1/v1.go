package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tianhongw/tinyid/handler/option"
	"github.com/tianhongw/tinyid/pkg/log"
	"github.com/tianhongw/tinyid/service"
)

type HandlerV1 struct {
	Services *service.Service
	Logger   log.Logger

	TinyId *TinyIdHandler
}

func NewHandlerV1(services *service.Service, opts *option.Options) *HandlerV1 {
	handler := &HandlerV1{
		Services: services,
		Logger:   opts.Logger,
	}

	handler.TinyId = NewTinyIdHandler(handler)

	return handler
}

func (h *HandlerV1) ResponseWithError(c *gin.Context, err error) {
	if err == nil {
		h.Logger.Warning("nil error")
		return
	}

	c.Error(err).SetType(gin.ErrorTypePublic)

	c.Abort()
}
