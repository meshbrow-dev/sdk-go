package meshbrow

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-key")
	if client.baseURL != "https://api.meshbrow.dev" {
		t.Errorf("expected default base URL, got %s", client.baseURL)
	}
	if client.apiKey != "test-key" {
		t.Errorf("expected api key test-key, got %s", client.apiKey)
	}
}

func TestWithBaseURL(t *testing.T) {
	client := NewClient("test-key", WithBaseURL("http://localhost:8080"))
	if client.baseURL != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %s", client.baseURL)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	client := NewClient("test-key", WithHTTPClient(custom))
	if client.httpClient != custom {
		t.Error("expected custom http client")
	}
}

func TestEnvelopeUnwrapping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"data":{"id":"mb_1","status":"ready"}}`)
	}))
	defer srv.Close()

	client := NewClient("key", WithBaseURL(srv.URL))
	sess, err := client.GetSession(context.Background(), "mb_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess.ID != "mb_1" {
		t.Errorf("expected ID mb_1, got %s", sess.ID)
	}
}

func TestNoEnvelopeFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// List sessions endpoint doesn't use data envelope
		fmt.Fprintf(w, `{"sessions":[{"id":"mb_1","status":"ready"}],"metrics":{"active_sessions":1}}`)
	}))
	defer srv.Close()

	client := NewClient("key", WithBaseURL(srv.URL))
	resp, err := client.ListSessions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(resp.Sessions))
	}
}

func TestEmptyBodyHandling(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient("key", WithBaseURL(srv.URL))
	err := client.DestroySession(context.Background(), "mb_1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error":{"code":"unauthorized","message":"Invalid API key"}}`)
	}))
	defer srv.Close()

	client := NewClient("bad-key", WithBaseURL(srv.URL))
	_, err := client.GetSession(context.Background(), "mb_1")
	if err == nil {
		t.Fatal("expected error")
	}

	// Should contain status code info
	errStr := err.Error()
	if errStr == "" {
		t.Error("expected non-empty error message")
	}
}

func TestUserAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "meshbrow-go/0.1.0" {
			t.Errorf("unexpected user agent: %s", r.Header.Get("User-Agent"))
		}
		fmt.Fprintf(w, `{"data":{"id":"mb_1","status":"ready"}}`)
	}))
	defer srv.Close()

	client := NewClient("key", WithBaseURL(srv.URL))
	_, _ = client.GetSession(context.Background(), "mb_1")
}
