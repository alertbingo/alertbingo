package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAlertLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"ok", 0, false},
		{"OK", 0, false},
		{"warn", 1, false},
		{"WARN", 1, false},
		{"alert", 2, false},
		{"ALERT", 2, false},
		{"invalid", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseAlertLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAlertLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseAlertLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatResponseSummary(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "all OK",
			input:    []byte(`[{"status":"OK"},{"status":"ok"}]`),
			expected: "",
		},
		{
			name:     "with errors",
			input:    []byte(`[{"status":"ERROR","errors":["bad request"]}]`),
			expected: "Statuses: ERROR; Errors: bad request",
		},
		{
			name:     "invalid JSON",
			input:    []byte(`not json`),
			expected: "not json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatResponseSummary(tt.input)
			if got != tt.expected {
				t.Errorf("FormatResponseSummary() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestClientSendChecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("expected Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected Content-Type header")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Response{{Status: "OK"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	checks := []CheckPayload{{
		Dashboard:  "test-dash",
		Site:       "test-site",
		Service:    "test-svc",
		Name:       "test-check",
		AlertLevel: 0,
	}}

	responses, err := client.SendChecks(context.Background(), checks)
	if err != nil {
		t.Fatalf("SendChecks() error = %v", err)
	}

	if len(responses) != 1 || responses[0].Status != "OK" {
		t.Errorf("unexpected response: %v", responses)
	}
}

func TestClientSendChecksError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`[{"status":"ERROR","errors":["invalid"]}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.SendChecks(context.Background(), []CheckPayload{})

	if err == nil {
		t.Error("expected error for bad status code")
	}
}
