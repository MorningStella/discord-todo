package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client represents an API client for Reclaim.ai
type ReclaimClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new Reclaim API client
func NewReclaimClient(baseURL string, token string) *ReclaimClient {
	return &ReclaimClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		token:      token,
	}
}

// ReclaimTask represents a task in Reclaim.ai
type ReclaimTask struct {
	Title              *string       `json:"title"`
	EventCategory      EventCategory `json:"eventCategory"`
	TimeChunksRequired uint8         `json:"timeChunksRequired"`
	MinChunkSize       *uint8        `json:"minChunkSize"`
	MaxChunkSize       *uint8        `json:"maxChunkSize"`
	AlwaysPrivate      bool          `json:"alwaysPrivate"`
	// TimeSchemeID       string        `json:"timeSchemeId"`
	Priority      TaskPriority `json:"priority"`
	ScheduleAfter *time.Time   `json:"snoozeUntil"`
	Due           *time.Time   `json:"due"`
	OnDeck        bool         `json:"onDeck"`
}

type TaskPriority string

const (
	PriorityP1 TaskPriority = "P1" // Highest priority
	PriorityP2 TaskPriority = "P2" // High priority
	PriorityP3 TaskPriority = "P3" // Medium priority
	PriorityP4 TaskPriority = "P4" // Low priority
)

type EventCategory string

const (
	Work     EventCategory = "WORK"
	Personal EventCategory = "PERSONAL"
)

// CreateTask creates a new task in Reclaim.ai
func (c *ReclaimClient) CreateTask(task ReclaimTask) (map[string]interface{}, error) {
	// Marshal the request body
	bodyData, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "https://api.app.reclaim.ai/api/tasks", bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response body
	log.Printf("Response body: %s", string(respBody))

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var responseObj map[string]interface{}
	if err := json.Unmarshal(respBody, &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return responseObj, nil
}
