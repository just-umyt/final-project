package http

import (
	"net/http"
	"time"
)

type ServerConfig struct {
	Address           string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
}

func NewServer(serverConfig *ServerConfig) *http.Server {
	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           serverConfig.Handler,
		ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
	}

	return server
}
