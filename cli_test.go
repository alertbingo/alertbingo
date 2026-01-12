package main

import (
	"testing"

	"github.com/alertbingo/cli/api"
)

func TestFormatResponses(t *testing.T) {
	tests := []struct {
		name      string
		responses []api.Response
		expected  string
	}{
		{
			name:      "all OK",
			responses: []api.Response{{Status: "OK"}, {Status: "ok"}},
			expected:  "",
		},
		{
			name:      "with non-OK status",
			responses: []api.Response{{Status: "ERROR"}},
			expected:  "Statuses: ERROR",
		},
		{
			name:      "with errors",
			responses: []api.Response{{Status: "OK", Errors: []string{"some error"}}},
			expected:  "Errors: some error",
		},
		{
			name:      "deduplicates statuses",
			responses: []api.Response{{Status: "Ignored"}, {Status: "Ignored"}, {Status: "Ignored"}},
			expected:  "Statuses: Ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatResponses(tt.responses)
			if got != tt.expected {
				t.Errorf("formatResponses() = %q, want %q", got, tt.expected)
			}
		})
	}
}
