package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Middleware func(next http.RoundTripper) http.RoundTripper

type loggingTransport struct {
	next   http.RoundTripper
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &loggingTransport{
			next:   next,
			logger: logger,
		}
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Duration("duration", duration),
	}

	if err != nil {
		t.logger.Error("API request failed", append(fields, zap.Error(err))...)
		return nil, err
	}

	if resp != nil {
		fields = append(fields, zap.Int("status", resp.StatusCode))

		if resp.StatusCode >= 500 {
			t.logger.Error("API server error", fields...)
		} else if resp.StatusCode >= 400 {
			t.logger.Warn("API client error", fields...)
		} else {
			t.logger.Debug("API request completed", fields...)
		}
	}

	return resp, nil
}

func ChainMiddleware(transport http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	for i := len(middlewares) - 1; i >= 0; i-- {
		transport = middlewares[i](transport)
	}
	return transport
}
