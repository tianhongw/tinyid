package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tianhongw/tinyid/handler"
	"github.com/tianhongw/tinyid/server/middleware/error_reporter"
	"github.com/tianhongw/tinyid/server/middleware/token"
)

func newRouter(handlers *handler.Handler) http.Handler {
	r := gin.Default()
	r.Use(error_reporter.JSONErrorReporter())

	handerV1 := handlers.V1

	apiV1 := r.Group("/tinyid", token.Authentication(handlers.V1.Services.Token))
	{
		apiV1.GET("/next_id", handerV1.TinyId.NextId)
	}

	return r
}
