package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var (
	clientReqID uint64
)

func getRequestID() string {
	myID := atomic.AddUint64(&clientReqID, 1)
	return fmt.Sprintf("%06d", myID)
}

// Transport implements http.RoundTripper. When set as Transport of http.Client, it executes HTTP requests with logging.
// No field is mandatory.
type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(resp *http.Response)
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
	ctx := req.Context()
	fields := []zap.Field{
		zap.String("host", req.Host),
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
	}
	if cv, ok := ctx.Value(ContextKeyRequestStart).(*contextValue); ok {
		fields = append(fields, zap.String("client_req_id", cv.reqID))
	}
	Log(req).Info("Outgoing HTTP Request", fields...)
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	logger := Log(resp.Request)
	fields := []zap.Field{
		zap.Int("status", resp.StatusCode),
		zap.String("uri", resp.Request.URL.String()),
		zap.String("host", resp.Request.Host),
	}
	if cv, ok := ctx.Value(ContextKeyRequestStart).(*contextValue); ok {
		elapsed := float64(time.Since(cv.start).Nanoseconds()) / 1000000.0
		fields = append(fields, zap.Float64("elapsed_ms", elapsed), zap.String("client_req_id", cv.reqID))
	}
	logger.Info("Incoming HTTP Response", fields...)
}

type contextKey struct {
	name string
}

type contextValue struct {
	start time.Time
	reqID string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ContextKeyRequestStart, &contextValue{
		start: time.Now(),
		reqID: getRequestID(),
	})
	req = req.WithContext(ctx)

	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(resp *http.Response) {
	if t.LogResponse != nil {
		t.LogResponse(resp)
	} else {
		DefaultLogResponse(resp)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
