package httpserver

import (
	"fmt"
	"time"
)

type Option func(*Server)

func Address(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf("%s:%s", host, port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
