package hoststats

import (
	"context"
	"testing"
)

func TestAppendAlertReason(t *testing.T) {
	tests := []struct {
		message  string
		reason   string
		expected string
	}{
		{"", "reason", "reason"},
		{"existing", "reason", "existing - reason"},
	}

	for _, tt := range tests {
		got := appendAlertReason(tt.message, tt.reason)
		if got != tt.expected {
			t.Errorf("appendAlertReason(%q, %q) = %q, want %q", tt.message, tt.reason, got, tt.expected)
		}
	}
}

func TestCollect(t *testing.T) {
	cfg := Config{
		Dashboard: "test-dash",
		Site:      "test-site",
		Service:   "test-svc",
	}

	checks, err := Collect(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Should return 5 checks: Memory, Uptime, CPU, Disk Used, Disk Inodes
	if len(checks) != 5 {
		t.Errorf("Collect() returned %d checks, want 5", len(checks))
	}

	expectedNames := []string{"Memory", "Uptime", "CPU", "Disk Used", "Disk Inodes"}
	for i, name := range expectedNames {
		if checks[i].Name != name {
			t.Errorf("checks[%d].Name = %q, want %q", i, checks[i].Name, name)
		}
		if checks[i].Dashboard != cfg.Dashboard {
			t.Errorf("checks[%d].Dashboard = %q, want %q", i, checks[i].Dashboard, cfg.Dashboard)
		}
	}
}
