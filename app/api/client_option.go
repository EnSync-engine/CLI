package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type ClientOption func(*Client)

func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Client) {
		c.log = logger
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.http.Timeout = timeout
	}
}

func WithRateLimit(rps float64, burst int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Limit(rps), burst)
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.http = httpClient
	}
}
