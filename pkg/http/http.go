package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type option func(*options)

type options struct {
	port    int
	handler http.Handler
}

func defaultOptions() *options {
	return &options{
		port: 8080,
	}
}

func (o *options) apply(opts ...option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPort(port int) option {
	return func(o *options) {
		o.port = port
	}
}

func WithHandler(handler http.Handler) option {
	return func(o *options) {
		o.handler = handler
	}
}

type HttpServer struct {
	opts   *options
	server *http.Server
}

func NewHttpServer(opts ...option) *HttpServer {
	o := defaultOptions()
	o.apply(opts...)
	return &HttpServer{
		opts: o,
	}
}

func (s *HttpServer) Start() error {
	s.server = &http.Server{
		Addr:    s.Addr(),
		Handler: s.opts.handler,
	}
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) Scheme() string {
	return "http"
}

func (s *HttpServer) Addr() string {
	return fmt.Sprintf(":%d", s.opts.port)
}
