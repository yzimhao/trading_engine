package http

import (
	"fmt"
	"net/http"
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
	opts *options
}

func NewHttpServer(opts ...option) *HttpServer {
	o := defaultOptions()
	o.apply(opts...)
	return &HttpServer{
		opts: o,
	}
}

func (s *HttpServer) Start() error {
	return http.ListenAndServe(s.Addr(), s.opts.handler)
}

func (s *HttpServer) Stop() error {
	return nil
}

func (s *HttpServer) Scheme() string {
	return "http"
}

func (s *HttpServer) Addr() string {
	return fmt.Sprintf(":%d", s.opts.port)
}
