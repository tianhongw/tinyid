package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tianhongw/tinyid/handler"
	"github.com/tianhongw/tinyid/pkg/log"
)

type Server struct {
	*http.Server

	logger log.Logger
}

type Option func(*Server)

func WithLogger(logger log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func NewServer(handlers *handler.Handler, port int, opts ...Option) (*Server, error) {
	if handlers == nil {
		return nil, errors.New("nil handlers")
	}

	s := &Server{
		logger: log.DummyLogger,
	}

	for _, o := range opts {
		o(s)
	}

	s.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: newRouter(handlers),
	}

	return s, nil
}

func (s *Server) Start() error {
	go func() {
		quit := make(chan os.Signal, 1)

		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		<-quit

		s.logger.Infoln("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Server.Shutdown(ctx); err != nil {
			s.logger.Warningf("server forced to shutdown: %v", err)
		}
	}()

	if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	s.logger.Infoln("server stopped.")
	return nil
}
