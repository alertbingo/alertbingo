// Package urlcheck provides functions to check URL availability and content.
package urlcheck

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alertbingo/alertbingo/api"
)

// Config holds the common configuration for URL checks
type Config struct {
	Dashboard        string
	Site             string
	Name             string
	Message          string
	InactiveExpire   string
	InactiveEscalate string
	Highlighted      string
	Timeout          time.Duration
}

// CheckParams holds the parameters for a single URL check
type CheckParams struct {
	URL          string
	ExpectedCode int    // 0 means any 2xx is OK
	ExpectedBody string // empty means no body check
}

// Result holds the result of a URL check
type Result struct {
	URL        string
	StatusCode int
	BodyMatch  bool
	Error      error
	Duration   time.Duration
	IsTimeout  bool
}

// Check performs a URL check and returns a CheckPayload
func Check(ctx context.Context, cfg Config, params CheckParams) api.CheckPayload {
	result := fetchURL(ctx, params, cfg.Timeout)
	return buildPayload(cfg, params, result)
}

// fetchURL performs the HTTP request and checks the response
func fetchURL(ctx context.Context, params CheckParams, timeout time.Duration) Result {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, params.URL, nil)
	if err != nil {
		return Result{URL: params.URL, Error: fmt.Errorf("invalid URL: %w", err)}
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		isTimeout := ctx.Err() == context.DeadlineExceeded ||
			strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "deadline exceeded")
		return Result{
			URL:       params.URL,
			Error:     fmt.Errorf("request failed: %w", err),
			Duration:  duration,
			IsTimeout: isTimeout,
		}
	}
	defer resp.Body.Close()

	result := Result{
		URL:        params.URL,
		StatusCode: resp.StatusCode,
		BodyMatch:  true,
		Duration:   duration,
	}

	// Check body content if required
	if params.ExpectedBody != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.Error = fmt.Errorf("failed to read body: %w", err)
			return result
		}
		result.BodyMatch = strings.Contains(string(body), params.ExpectedBody)
	}

	return result
}

// slowThreshold is the duration above which a successful request is considered slow
const slowThreshold = 2 * time.Second

// buildPayload creates a CheckPayload from the check result
func buildPayload(cfg Config, params CheckParams, result Result) api.CheckPayload {
	payload := api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          params.URL,
		Name:             cfg.Name,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}

	durationSecs := result.Duration.Seconds()

	// Handle request errors
	if result.Error != nil {
		payload.AlertLevel = 2 // alert
		if result.IsTimeout {
			payload.Value = "Timeout"
		} else {
			payload.Value = "Error"
		}
		payload.Message = appendAlertReason(cfg.Message, result.Error.Error())
		return payload
	}

	// Check status code
	codeOK := false
	if params.ExpectedCode == 0 {
		// Any 2xx is OK
		codeOK = result.StatusCode >= 200 && result.StatusCode < 300
	} else {
		codeOK = result.StatusCode == params.ExpectedCode
	}

	// Determine alert level
	var reasons []string

	if !codeOK {
		if params.ExpectedCode == 0 {
			reasons = append(reasons, fmt.Sprintf("expected status 2xx"))
		} else {
			reasons = append(reasons, fmt.Sprintf("expected status %d", params.ExpectedCode))
		}
	}

	if !result.BodyMatch {
		reasons = append(reasons, fmt.Sprintf("body missing expected string: %q", params.ExpectedBody))
	}

	if len(reasons) > 0 {
		payload.AlertLevel = 2 // alert
		payload.Value = fmt.Sprintf("%d", result.StatusCode)
		payload.Message = appendAlertReason(cfg.Message, strings.Join(reasons, "; "))
	} else if result.Duration > slowThreshold {
		// Successful but slow response - warning level
		payload.AlertLevel = 1 // warning
		payload.Value = fmt.Sprintf("%d", result.StatusCode)
		payload.Message = appendAlertReason(cfg.Message, fmt.Sprintf("slow response: %.2fs", durationSecs))
	} else {
		payload.AlertLevel = 0 // ok
		payload.Value = fmt.Sprintf("%d", result.StatusCode)
		payload.Message = cfg.Message
	}

	return payload
}

// appendAlertReason appends an alert reason to an existing message
func appendAlertReason(message, reason string) string {
	if message == "" {
		return reason
	}
	return message + " - " + reason
}
