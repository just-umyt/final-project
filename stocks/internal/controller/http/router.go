package http

import (
	"net/http"
	"time"
)

type ServerConfig struct {
	Addr              string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
}

func NewServer(serverConfig *ServerConfig) *http.Server {
	server := &http.Server{
		Addr:              serverConfig.Addr,
		Handler:           serverConfig.Handler,
		ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
	}

	return server
}
