package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServe(address string, timeout time.Duration) error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:        address,
		Handler:     mux,
		ReadTimeout: timeout,
	}

	return server.ListenAndServe()
}
