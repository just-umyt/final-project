package http

import (
	"net/http"
	"time"
)

type SeverConfig struct {
	Addr              string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
}

func NewServer(serverConfig *SeverConfig) *http.Server {

	server := &http.Server{
		Addr:              serverConfig.Addr,
		ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
	}

	return server
}
