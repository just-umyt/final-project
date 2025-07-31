package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics interface {
	IncRequest(path string)
	ObserveLatency(path string, duration float64)
	IncError(path string)
}

var _ Metrics = &AppMetrics{}

type AppMetrics struct {
	RequestCounter  *prometheus.CounterVec
	ResponseLatency *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
}

func (a *AppMetrics) IncRequest(path string) {
	a.RequestCounter.With(prometheus.Labels{"path": path}).Inc()
}

func (a *AppMetrics) ObserveLatency(path string, duration float64) {
	a.ResponseLatency.With(prometheus.Labels{"path": path}).Observe(duration)
}

func (a *AppMetrics) IncError(path string) {
	a.ErrorsTotal.With(prometheus.Labels{"path": path}).Inc()
}

func RegisterMetrics() *AppMetrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)

	responseLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Histogram of response latencies",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_failed_requests_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"path"},
	)

	prometheus.MustRegister(requestCounter, responseLatency, errorCounter)

	return &AppMetrics{
		RequestCounter:  requestCounter,
		ResponseLatency: responseLatency,
		ErrorsTotal:     errorCounter,
	}
}
