package urlcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	}))
	defer server.Close()

	cfg := Config{
		Dashboard: "test",
		Site:      "test-site",
		Timeout:   5 * time.Second,
	}

	params := CheckParams{
		URL: server.URL,
	}

	result := Check(context.Background(), cfg, params)

	if result.AlertLevel != 0 {
		t.Errorf("expected AlertLevel 0, got %d", result.AlertLevel)
	}
	if result.Value != "200" {
		t.Errorf("expected Value '200', got %s", result.Value)
	}
}

func TestCheck_ExpectedCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	cfg := Config{
		Dashboard: "test",
		Site:      "test-site",
		Timeout:   5 * time.Second,
	}

	// Should pass with expected code 201
	params := CheckParams{
		URL:          server.URL,
		ExpectedCode: 201,
	}

	result := Check(context.Background(), cfg, params)
	if result.AlertLevel != 0 {
		t.Errorf("expected AlertLevel 0, got %d", result.AlertLevel)
	}

	// Should fail without expected code (201 is 2xx so should pass)
	params2 := CheckParams{
		URL: server.URL,
	}

	result2 := Check(context.Background(), cfg, params2)
	if result2.AlertLevel != 0 {
		t.Errorf("expected AlertLevel 0 for 2xx, got %d", result2.AlertLevel)
	}
}

func TestCheck_WrongCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := Config{
		Dashboard: "test",
		Site:      "test-site",
		Timeout:   5 * time.Second,
	}

	params := CheckParams{
		URL: server.URL,
	}

	result := Check(context.Background(), cfg, params)

	if result.AlertLevel != 2 {
		t.Errorf("expected AlertLevel 2 for 404, got %d", result.AlertLevel)
	}
}

func TestCheck_BodyMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	}))
	defer server.Close()

	cfg := Config{
		Dashboard: "test",
		Site:      "test-site",
		Timeout:   5 * time.Second,
	}

	// Should pass with matching body
	params := CheckParams{
		URL:          server.URL,
		ExpectedBody: "Hello",
	}

	result := Check(context.Background(), cfg, params)
	if result.AlertLevel != 0 {
		t.Errorf("expected AlertLevel 0 for matching body, got %d", result.AlertLevel)
	}

	// Should fail with non-matching body
	params2 := CheckParams{
		URL:          server.URL,
		ExpectedBody: "Goodbye",
	}

	result2 := Check(context.Background(), cfg, params2)
	if result2.AlertLevel != 2 {
		t.Errorf("expected AlertLevel 2 for non-matching body, got %d", result2.AlertLevel)
	}
}

func TestCheck_ConnectionError(t *testing.T) {
	cfg := Config{
		Dashboard: "test",
		Site:      "test-site",
		Timeout:   1 * time.Second,
	}

	params := CheckParams{
		URL: "http://localhost:99999",
	}

	result := Check(context.Background(), cfg, params)

	if result.AlertLevel != 2 {
		t.Errorf("expected AlertLevel 2 for connection error, got %d", result.AlertLevel)
	}
	if result.Value != "Error" {
		t.Errorf("expected Value 'Error', got %s", result.Value)
	}
}
