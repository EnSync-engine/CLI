package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func buildURL(baseURL, path string, query url.Values) string {
	reqURL := baseURL + path
	if len(query) > 0 {
		reqURL = fmt.Sprintf("%s?%s", reqURL, query.Encode())
	}
	return reqURL
}

func marshalBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return bytes.NewReader(bodyBytes), nil
}

func buildRequest(ctx context.Context, method, reqURL string, body io.Reader, accessKey string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(headerAccessKey, accessKey)
	req.Header.Set(headerAccept, contentTypeJSON)
	if body != nil {
		req.Header.Set(headerContentType, contentTypeJSON)
	}

	return req, nil
}
