package log

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/sirupsen/logrus"
)

type HttpLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *HttpLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func (l *HttpLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func NewMiddlewareLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&MiddlewareLogger{httpLogger})
}

type MiddlewareLogger struct {
	Logger *logrus.Logger
}

func (l *MiddlewareLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &HttpLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC3339)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Infoln("request started")

	return entry
}

func GetHttpLogger(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*HttpLoggerEntry)
	return entry.Logger
}
