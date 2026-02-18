package certcheck

import (
	"context"
	"testing"
	"time"
)

func TestCollect_ValidCertificate(t *testing.T) {
	cfg := Config{
		Dashboard: "test-dash",
		Site:      "test-site",
		Name:      "cert-check",
		Timeout:   10 * time.Second,
	}

	// Test with a well-known site that should have a valid certificate
	checks := Collect(context.Background(), cfg, []string{"https://google.com"})

	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}

	check := checks[0]
	if check.Dashboard != "test-dash" {
		t.Errorf("expected dashboard 'test-dash', got %s", check.Dashboard)
	}
	if check.Site != "test-site" {
		t.Errorf("expected site 'test-site', got %s", check.Site)
	}
	// Service is set to the host from the certificate
	if check.Service != "google.com" {
		t.Errorf("expected service 'google.com', got %s", check.Service)
	}
	if check.Name != "cert-check" {
		t.Errorf("expected name 'cert-check', got %s", check.Name)
	}
	// Google's cert should be valid for more than 2 weeks
	if check.AlertLevel != 0 {
		t.Errorf("expected alert level 0, got %d", check.AlertLevel)
	}
}

func TestCollect_InvalidURL(t *testing.T) {
	cfg := Config{
		Dashboard: "test-dash",
		Site:      "test-site",
		Name:      "cert-check",
		Timeout:   5 * time.Second,
	}

	// Test with an invalid host
	checks := Collect(context.Background(), cfg, []string{"https://invalid.invalid.invalid"})

	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}

	check := checks[0]
	// Should be alert level 2 for error
	if check.AlertLevel != 2 {
		t.Errorf("expected alert level 2 for error, got %d", check.AlertLevel)
	}
	if check.Value != "Error" {
		t.Errorf("expected value 'Error', got %s", check.Value)
	}
	if check.Message == "" {
		t.Error("expected error message to be set")
	}
}

func TestCollect_MultipleURLs(t *testing.T) {
	cfg := Config{
		Dashboard: "test-dash",
		Site:      "test-site",
		Name:      "cert-check",
		Timeout:   10 * time.Second,
	}

	urls := []string{"https://google.com", "https://github.com"}
	checks := Collect(context.Background(), cfg, urls)

	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}
}

func TestAppendAlertReason(t *testing.T) {
	tests := []struct {
		message  string
		reason   string
		expected string
	}{
		{"", "test reason", "test reason"},
		{"existing", "new reason", "existing - new reason"},
	}

	for _, tt := range tests {
		result := appendAlertReason(tt.message, tt.reason)
		if result != tt.expected {
			t.Errorf("appendAlertReason(%q, %q) = %q, want %q", tt.message, tt.reason, result, tt.expected)
		}
	}
}
