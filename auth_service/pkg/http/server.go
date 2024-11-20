package http

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	addr   string
	mux    http.Handler
}

func NewServer(addr string, handler http.Handler) *Server {
	s := &Server{addr: addr, mux: handler}
	s.server = &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
