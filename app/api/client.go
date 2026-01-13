package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/EnSync-engine/CLI/app/domain"
)

var _ APIClient = (*Client)(nil)

const (
	defaultRetryMax     = 3
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 5 * time.Second
)

type Client struct {
	baseURL   string
	accessKey string

	http        *http.Client
	rateLimiter *rate.Limiter
	log         *zap.Logger
}

func NewClient(baseURL string, options ...ClientOption) *Client {
	retryable := retryablehttp.NewClient()
	retryable.RetryMax = defaultRetryMax
	retryable.RetryWaitMin = defaultRetryWaitMin
	retryable.RetryWaitMax = defaultRetryWaitMax
	retryable.Logger = nil

	client := &Client{
		baseURL: baseURL,
		http:    retryable.StandardClient(),
		log:     zap.NewNop(),
	}

	for _, opt := range options {
		opt(client)
	}

	if client.log != zap.NewNop() {
		client.http.Transport = NewLoggingMiddleware(client.log)(client.http.Transport)
	}

	return client
}

func (c *Client) SetAccessKey(key string) {
	c.accessKey = key
}

func (c *Client) execute(ctx context.Context, method, path string, queryParams url.Values, requestBody any) ([]byte, error) {
	if err := c.waitForRateLimit(ctx); err != nil {
		return nil, err
	}

	fullURL := buildURL(c.baseURL, path, queryParams)
	body, err := marshalBody(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := buildRequest(ctx, method, fullURL, body, c.accessKey)
	if err != nil {
		return nil, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := readResponseBody(response)
	if err != nil {
		return nil, err
	}

	if isErrorStatus(response.StatusCode) {
		return nil, handleErrorResponse(response.StatusCode, responseBody)
	}

	return responseBody, nil
}

func (c *Client) waitForRateLimit(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}
	return nil
}

func (c *Client) ListEvents(ctx context.Context, params *ListParams) (*domain.EventList, error) {
	responseData, err := c.execute(ctx, http.MethodGet, pathEvent, params.ToQuery(), nil)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	var events domain.EventList
	if err := unmarshalResponse(responseData, &events); err != nil {
		return nil, err
	}

	return &events, nil
}

func (c *Client) GetEventByName(ctx context.Context, eventName string) (*domain.Event, error) {
	path := fmt.Sprintf(pathEventByName, url.PathEscape(eventName))

	responseData, err := c.execute(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get event %q: %w", eventName, err)
	}

	var event domain.Event
	if err := unmarshalResponse(responseData, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Client) CreateEvent(ctx context.Context, event *domain.Event) error {
	if _, err := c.execute(ctx, http.MethodPost, pathEvent, nil, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (c *Client) UpdateEvent(ctx context.Context, event *domain.Event) error {
	path := fmt.Sprintf(pathEventByID, event.ID)
	updatePayload := map[string]any{
		"name":    event.Name,
		"payload": event.Payload,
	}

	if _, err := c.execute(ctx, http.MethodPut, path, nil, updatePayload); err != nil {
		return fmt.Errorf("update event (id=%s): %w", event.ID, err)
	}
	return nil
}

func (c *Client) ListAccessKeys(ctx context.Context, params *ListParams) (*domain.AccessKeyList, error) {
	responseData, err := c.execute(ctx, http.MethodGet, pathAccessKey, params.ToQuery(), nil)
	if err != nil {
		return nil, fmt.Errorf("list access keys: %w", err)
	}

	var accessKeys domain.AccessKeyList
	if err := unmarshalResponse(responseData, &accessKeys); err != nil {
		return nil, err
	}

	return &accessKeys, nil
}

func (c *Client) GetAccessKeyByID(ctx context.Context, id string) (*domain.AccessKeyPermissions, error) {
	path := fmt.Sprintf(pathAccessKeyByID, url.PathEscape(id))

	responseData, err := c.execute(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get access key %q: %w", id, err)
	}

	var key domain.AccessKeyPermissions
	if err := unmarshalResponse(responseData, &key); err != nil {
		return nil, err
	}

	return &key, nil
}

func (c *Client) CreateAccessKey(ctx context.Context, req *domain.CreateAccessKeyRequest) (*domain.AccessKey, error) {
	responseData, err := c.execute(ctx, http.MethodPost, pathAccessKey, nil, req)
	if err != nil {
		return nil, fmt.Errorf("create access key: %w", err)
	}

	var key domain.AccessKey
	if err := unmarshalResponse(responseData, &key); err != nil {
		return nil, err
	}

	return &key, nil
}

func (c *Client) DeleteAccessKey(ctx context.Context, id string) error {
	path := fmt.Sprintf(pathAccessKeyByID, url.PathEscape(id))

	if _, err := c.execute(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("delete access key %q: %w", id, err)
	}
	return nil
}

func (c *Client) GetAccessKeyPermissions(ctx context.Context, accessKey string) (*domain.AccessKeyPermissions, error) {
	path := fmt.Sprintf(pathAccessKeyPermissions, url.PathEscape(accessKey))

	responseData, err := c.execute(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get permissions for key %q: %w", accessKey, err)
	}

	var permissions domain.AccessKeyPermissions
	if err := unmarshalResponse(responseData, &permissions); err != nil {
		return nil, err
	}

	return &permissions, nil
}

func (c *Client) SetAccessKeyPermissions(ctx context.Context, accessKey string, permissions *domain.Permissions) error {
	path := fmt.Sprintf(pathAccessKeyPermissions, url.PathEscape(accessKey))
	updatePayload := map[string]any{
		"send":    permissions.Send,
		"receive": permissions.Receive,
	}

	if _, err := c.execute(ctx, http.MethodPost, path, nil, updatePayload); err != nil {
		return fmt.Errorf("set permissions for key %q: %w", accessKey, err)
	}
	return nil
}

func (c *Client) UpdateServiceKeyPair(ctx context.Context, accessKey string) (*domain.ServiceKeyPair, error) {
	req := &domain.UpdateServiceKeyPairRequest{AccessKey: accessKey}

	responseData, err := c.execute(ctx, http.MethodPut, pathServiceKeyPair, nil, req)
	if err != nil {
		return nil, fmt.Errorf("update service key pair: %w", err)
	}

	var keyPair domain.ServiceKeyPair
	if err := unmarshalResponse(responseData, &keyPair); err != nil {
		return nil, err
	}

	return &keyPair, nil
}

func (c *Client) ListWorkspaces(ctx context.Context, params *ListParams) (*domain.WorkspaceList, error) {
	responseData, err := c.execute(ctx, http.MethodGet, pathWorkspace, params.ToQuery(), nil)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}

	var workspaces domain.WorkspaceList
	if err := unmarshalResponse(responseData, &workspaces); err != nil {
		return nil, err
	}

	return &workspaces, nil
}

func (c *Client) CreateWorkspace(ctx context.Context, name string) error {
	req := &domain.CreateWorkspaceRequest{Name: name}

	if _, err := c.execute(ctx, http.MethodPost, pathWorkspace, nil, req); err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	return nil
}
