package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("invalid response: nil response or body")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func handleErrorResponse(statusCode int, body []byte) error {
	return fmt.Errorf("request failed with status %d: %s", statusCode, string(body))
}

func unmarshalResponse[T any](data []byte, target *T) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

func isErrorStatus(statusCode int) bool {
	return statusCode >= http.StatusBadRequest
}
