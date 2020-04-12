package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"

	"gitlab.unanet.io/devops/eve/pkg/metrics"
)

// Metrics adapts the incoming request with Logging/Metrics
func Metrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Tally the incoming request metrics
		now := time.Now()
		metrics.StatHTTPRequestCount.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()
		metrics.StatRequestSaturationGauge.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Run this on the way out (i.e. outgoing response)
		defer func() {
			// Calculate the request duration (i.e. latency)
			ms := float64(time.Since(now).Nanoseconds()) / 1000000.0

			// Tally the outgoing response metrics
			metrics.StatRequestDurationHistogram.WithLabelValues(r.RequestURI, r.Method, r.Proto).Observe(ms)
			metrics.StatRequestDurationGauge.WithLabelValues(r.RequestURI, r.Method, r.Proto).Set(ms)
			metrics.StatRequestSaturationGauge.WithLabelValues(r.RequestURI, r.Method, r.Proto).Dec()
			metrics.StatHTTPResponseCount.WithLabelValues(strconv.Itoa(ww.Status()), r.RequestURI, r.Method, r.Proto).Inc()
		}()
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
