// Package api provides a client for the Alert Bingo API.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CheckPayload represents a single check to send to the API
type CheckPayload struct {
	Dashboard        string `json:"dashboard"`
	Site             string `json:"site"`
	Service          string `json:"service"`
	Name             string `json:"name"`
	AlertLevel       int    `json:"alert_level"`
	Message          string `json:"message,omitempty"`
	Value            string `json:"value,omitempty"`
	InactiveExpire   string `json:"inactive_expire,omitempty"`
	InactiveEscalate string `json:"inactive_escalate,omitempty"`
	Highlighted      string `json:"highlighted,omitempty"`
}

// Response represents the response from the API
type Response struct {
	Status string   `json:"status"`
	Errors []string `json:"errors,omitempty"`
}

// Client is an API client for Alert Bingo
type Client struct {
	APIURL     string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(apiURL, token string) *Client {
	return &Client{
		APIURL:     apiURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// SendChecks sends a slice of checks to the API
func (c *Client) SendChecks(ctx context.Context, checks []CheckPayload) ([]Response, error) {
	jsonData, err := json.Marshal(checks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, FormatResponseSummary(body))
	}

	var responses []Response
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, truncateBody(body, 200))
	}

	return responses, nil
}

// ParseAlertLevel converts a string alert level to an integer
func ParseAlertLevel(level string) (int, error) {
	switch strings.ToLower(level) {
	case "ok":
		return 0, nil
	case "warn":
		return 1, nil
	case "alert":
		return 2, nil
	default:
		return -1, fmt.Errorf("invalid alert level: %s (must be ok, warn, or alert)", level)
	}
}

// FormatResponseSummary returns a formatted summary of non-OK statuses and errors from API response array
func FormatResponseSummary(body []byte) string {
	var responses []Response
	if err := json.Unmarshal(body, &responses); err != nil {
		// If JSON parsing fails, return the raw body
		return string(body)
	}

	// Collect unique non-OK statuses
	statusSet := make(map[string]struct{})
	for _, resp := range responses {
		if strings.ToUpper(resp.Status) != "OK" {
			statusSet[resp.Status] = struct{}{}
		}
	}

	// Collect unique errors
	errorSet := make(map[string]struct{})
	for _, resp := range responses {
		for _, e := range resp.Errors {
			errorSet[e] = struct{}{}
		}
	}

	var parts []string

	// Format unique non-OK statuses
	if len(statusSet) > 0 {
		var statuses []string
		for s := range statusSet {
			statuses = append(statuses, s)
		}
		parts = append(parts, fmt.Sprintf("Statuses: %s", strings.Join(statuses, ", ")))
	}

	// Format unique errors
	if len(errorSet) > 0 {
		var errors []string
		for e := range errorSet {
			errors = append(errors, e)
		}
		parts = append(parts, fmt.Sprintf("Errors: %s", strings.Join(errors, ", ")))
	}

	return strings.Join(parts, "; ")
}

// truncateBody truncates a byte slice to a maximum length, appending "..." if truncated
func truncateBody(body []byte, maxLen int) string {
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "..."
}
