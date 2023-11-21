package util

import (
	"fmt"
	"io"
	"net/http"
)

func MakeAPICall(method, endpoint string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get response: %w", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed read response: %w", err)
	}

	return data, resp.StatusCode, nil
}
