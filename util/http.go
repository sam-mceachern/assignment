package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakeAPICall[O any](method, endpoint string) (O, error) {
	var responseBody O
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return responseBody, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return responseBody, fmt.Errorf("failed to get response: %w", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseBody, fmt.Errorf("failed read response: %w", err)
	}

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		return responseBody, fmt.Errorf("failed unmarhsal json: %w", err)
	}

	return responseBody, nil
}
