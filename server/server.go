package server

import (
	"net/http"

	"github.com/tianhongw/tinyid/pkg/log"
)

type Server struct {
	*http.Server

	log.Logger
}
