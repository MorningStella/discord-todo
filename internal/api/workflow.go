package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MorningStella/discord-todo/internal/config"
)

// WorkflowClient represents an API client for external services
type WorkflowClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewWorkflowClient creates a new API client
func NewWorkflowClient(cfg *config.Config) *WorkflowClient {
	return &WorkflowClient{
		baseURL:    cfg.WorkflowBaseUrl,
		httpClient: &http.Client{},
	}
}

func (c *WorkflowClient) SendTodoRequest(requestBody map[string]interface{}, action string) (string, error) {
	// Marshal the request body
	bodyData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request body: %w", err)
	}

	// Create and send request
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(bodyData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	query := req.URL.Query()
	query.Add("action", action)
	req.URL.RawQuery = query.Encode()

	// Create HTTP client and send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response body
	log.Printf("Response body: %s", string(respBody))

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("server responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response to get the message
	var responseObj map[string]interface{}
	if err := json.Unmarshal(respBody, &responseObj); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the message
	message, ok := responseObj["message"].(string)
	if !ok {
		return "", fmt.Errorf("response does not contain a message field")
	}

	return message, nil
}
