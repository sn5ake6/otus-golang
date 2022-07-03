package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	addr       string
	logger     Logger
	httpServer *http.Server
}

type Logger interface {
	Error(message string)
	Warning(msg string)
	Info(message string)
	Debug(msg string)
	LogRequest(r *http.Request, statusCode int, requestDuration time.Duration)
}

type Application interface { // TODO
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		addr:   addr,
		logger: logger,
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(http.HandlerFunc(s.helloWorld), logger),
	}

	s.httpServer = httpServer

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("Server started: %s", s.addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("Server stopped: %s", s.addr))

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello-world"))
}
