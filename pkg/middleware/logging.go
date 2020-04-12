package middleware

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.unanet.io/devops/eve/pkg/log"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type LogEntry struct {
	logger *zap.Logger
}

func (l *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.logger.Info("Outgoing HTTP Response",
		zap.Int("http_status", status),
		zap.Int("resp_bytes_length", bytes),
		zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0))
}

func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.logger.Panic(fmt.Sprintf("%+v", v),
		zap.String("stack", string(stack)))
}

func Logger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&LogEntryConstructor{log.Logger})
}

type LogEntryConstructor struct {
	logger *zap.Logger
}

func (l *LogEntryConstructor) NewLogEntry(r *http.Request) middleware.LogEntry {
	var logFields []zap.Field

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, zap.String("req_id", reqID))
	}

	incomingRequestFields := []zap.Field{
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	}

	entry := &LogEntry{
		logger: l.logger.With(logFields...),
	}

	l.logger.With(incomingRequestFields...).Info("Incoming HTTP Request")

	return entry
}

func Log(r *http.Request) *zap.Logger {
	entry := middleware.GetLogEntry(r).(*LogEntry)
	return entry.logger
}
