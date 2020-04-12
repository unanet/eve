package metrics

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort = 3000
)

var (
	StatRequestSaturationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_saturation",
			Help: "The total number of requests inside the server (transactions serving)",
		}, []string{"uri", "method", "protocol"})

	StatBuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_build_info",
			Help: "A metric with a constant '1' value labeled by version, revision, branch, and goversion from which the service was build was built.",
		}, []string{"service", "revision", "branch", "version", "author", "build_date", "build_user", "build_host"})

	StatRequestDurationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_duration_ms",
			Help: "The time the server spends processing a request in milliseconds",
		}, []string{"uri", "method", "protocol"})

	StatAuditCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_total",
			Help: "The total number of audit events",
		}, []string{"event"})

	StatHTTPRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number of incoming requests to the service",
		}, []string{"uri", "method", "protocol"})

	StatHTTPResponseCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_total",
			Help: "The total number of outgoing responses to the client",
		}, []string{"code", "uri", "method", "protocol"})

	StatRequestDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_histogram_ms",
			Help:    "time spent processing an http request in milliseconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 18),
		}, []string{"uri", "method", "protocol"})
)

func StartMetrics() {
	var c Config
	configErr := envconfig.Process("EVE", &c)
	if configErr != nil {
		c.PromPort = defaultPort
	}

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%v", c.PromPort), nil)
}
