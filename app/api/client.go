package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/EnSync-engine/CLI/app/domain"
)

// Client is the API client for interacting with the EnSync API.
type Client struct {
	baseURL     string
	accessKey   string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	logger      *zap.Logger
}

// NewClient creates a new API client with the given base URL and options.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second

	c := &Client{
		baseURL:    baseURL,
		httpClient: retryClient.StandardClient(),
		logger:     zap.NewNop(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetAccessKey sets the access key for API authentication.
func (c *Client) SetAccessKey(accessKey string) {
	c.accessKey = accessKey
}

// doRequest performs an HTTP request with the given method, path, query parameters, and body.
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, error) {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}
	}

	reqURL := c.baseURL + path
	if query != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, query.Encode())
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(XAccessKeyHeader, c.accessKey)
	if body != nil {
		req.Header.Set("Content-Type", ContentTypeHeader)
	}
	req.Header.Set("Accept", ContentTypeHeader)

	c.logger.Debug("Sending request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("url", reqURL),
		zap.String("query", query.Encode()),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.logger.Debug("Received response",
		zap.Int("status", resp.StatusCode),
		zap.Int("body_size", len(respBody)),
	)

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	return respBody, nil
}

// ListEvents retrieves a list of events based on the provided parameters.
func (c *Client) ListEvents(ctx context.Context, params *ListParams) (*domain.EventList, error) {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", params.PageIndex))
	query.Set("limit", fmt.Sprintf("%d", params.Limit))
	query.Set("order", params.Order)
	query.Set("orderBy", params.OrderBy)

	data, err := c.doRequest(ctx, http.MethodGet, "/event", query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var response domain.EventList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CreateEvent creates a new event.
func (c *Client) CreateEvent(ctx context.Context, event *domain.Event) error {
	_, err := c.doRequest(ctx, http.MethodPost, "/event", nil, event)
	return err
}

// UpdateEvent updates an existing event.
func (c *Client) UpdateEvent(ctx context.Context, event *domain.Event) error {
	url := fmt.Sprintf("/event/%d", event.ID)
	payload := map[string]interface{}{
		"name":    event.Name,
		"payload": event.Payload,
	}

	_, err := c.doRequest(ctx, http.MethodPut, url, nil, payload)
	return err
}

// GetEventByName retrieves an event by its name.
func (c *Client) GetEventByName(ctx context.Context, name string) (*domain.Event, error) {
	encodedName := url.PathEscape(name)
	url := fmt.Sprintf("/event/%s", encodedName)

	data, err := c.doRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by name '%s': %w", name, err)
	}

	var event domain.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event response for name '%s': %w", name, err)
	}

	return &event, nil
}

// ListAccessKeys retrieves a list of access keys.
func (c *Client) ListAccessKeys(ctx context.Context, params *ListParams) (*domain.AccessKeyList, error) {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", params.PageIndex))
	query.Set("limit", fmt.Sprintf("%d", params.Limit))
	query.Set("order", params.Order)
	query.Set("orderBy", params.OrderBy)

	data, err := c.doRequest(ctx, http.MethodGet, "/access-key", query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list access keys: %w", err)
	}

	var response domain.AccessKeyList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// CreateAccessKey creates a new access key.
func (c *Client) CreateAccessKey(ctx context.Context, permissions *domain.Permissions) (*domain.AccessKey, error) {
	payload := map[string]interface{}{
		"permissions": permissions,
	}

	data, err := c.doRequest(ctx, http.MethodPost, "/access-key", nil, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create access key: %w", err)
	}

	var response domain.AccessKey
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAccessKeyPermissions retrieves permissions for a specific access key.
func (c *Client) GetAccessKeyPermissions(ctx context.Context, key string) (*domain.AccessKeyPermissions, error) {
	encodedKey := url.PathEscape(key)
	url := fmt.Sprintf("/access-key/permissions/%s", encodedKey)

	data, err := c.doRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get access key permissions: %w", err)
	}

	var response domain.AccessKeyPermissions
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SetAccessKeyPermissions updates permissions for a specific access key.
func (c *Client) SetAccessKeyPermissions(ctx context.Context, key string, permissions *domain.Permissions) error {
	url := fmt.Sprintf("/access-key/permissions/%s", key)
	payload := map[string]interface{}{
		"send":    permissions.Send,
		"receive": permissions.Receive,
	}

	_, err := c.doRequest(ctx, http.MethodPost, url, nil, payload)
	return err
}
