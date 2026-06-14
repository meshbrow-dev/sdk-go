package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer creates a test HTTP server that responds with predefined responses.
func testServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := NewClient("test-api-key", WithBaseURL(srv.URL))
	return client, srv
}

func TestCreateSession(t *testing.T) {
	t.Run("basic create", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v1/sessions" {
				t.Errorf("expected /v1/sessions, got %s", r.URL.Path)
			}
			if r.Header.Get("Authorization") != "Bearer test-api-key" {
				t.Errorf("missing auth header")
			}

			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["stealth"] != "max" {
				t.Errorf("expected stealth=max, got %v", body["stealth"])
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"data":{"id":"mb_test1","status":"ready","cdp_endpoint":"wss://api.meshbrow.dev/cdp/mb_test1","token":"tok1","created_at":"2026-06-14T00:00:00Z","expires_at":"2026-06-14T01:00:00Z"}}`)
		})
		defer srv.Close()

		sess, err := client.CreateSession(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sess.ID != "mb_test1" {
			t.Errorf("expected ID mb_test1, got %s", sess.ID)
		}
		if sess.Status != "ready" {
			t.Errorf("expected status ready, got %s", sess.Status)
		}
		if sess.CDPEndpoint != "wss://api.meshbrow.dev/cdp/mb_test1" {
			t.Errorf("unexpected cdp_endpoint: %s", sess.CDPEndpoint)
		}
	})

	t.Run("with proxy", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			proxy, ok := body["proxy"].(map[string]any)
			if !ok {
				t.Fatal("expected proxy in body")
			}
			if proxy["type"] != "residential" {
				t.Errorf("expected proxy type residential, got %v", proxy["type"])
			}
			if proxy["country"] != "US" {
				t.Errorf("expected proxy country US, got %v", proxy["country"])
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"data":{"id":"mb_proxy1","status":"ready"}}`)
		})
		defer srv.Close()

		sess, err := client.CreateSession(context.Background(), &CreateSessionParams{
			Stealth:      "max",
			ProxyType:    "residential",
			ProxyCountry: "US",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sess.ID != "mb_proxy1" {
			t.Errorf("expected ID mb_proxy1, got %s", sess.ID)
		}
	})

	t.Run("with viewport", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			vp, ok := body["viewport"].(map[string]any)
			if !ok {
				t.Fatal("expected viewport in body")
			}
			if vp["width"] != float64(1920) {
				t.Errorf("expected width 1920, got %v", vp["width"])
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"data":{"id":"mb_vp1","status":"ready"}}`)
		})
		defer srv.Close()

		sess, err := client.CreateSession(context.Background(), &CreateSessionParams{
			Viewport: &Viewport{Width: 1920, Height: 1080},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sess.ID != "mb_vp1" {
			t.Errorf("expected ID mb_vp1, got %s", sess.ID)
		}
	})
}

func TestGetSession(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/sessions/mb_abc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		fmt.Fprintf(w, `{"data":{"id":"mb_abc123","status":"ready","stealth":"max","created_at":"2026-06-14T00:00:00Z"}}`)
	})
	defer srv.Close()

	sess, err := client.GetSession(context.Background(), "mb_abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess.ID != "mb_abc123" {
		t.Errorf("expected ID mb_abc123, got %s", sess.ID)
	}
	if sess.Stealth != "max" {
		t.Errorf("expected stealth max, got %s", sess.Stealth)
	}
}

func TestListSessions(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/sessions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// List endpoint returns without envelope (directly)
		fmt.Fprintf(w, `{"sessions":[{"id":"mb_1","status":"ready"},{"id":"mb_2","status":"ready"}],"metrics":{"active_sessions":2}}`)
	})
	defer srv.Close()

	resp, err := client.ListSessions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(resp.Sessions))
	}
	if resp.Sessions[0].ID != "mb_1" {
		t.Errorf("expected first session ID mb_1, got %s", resp.Sessions[0].ID)
	}
}

func TestDestroySession(t *testing.T) {
	t.Run("basic destroy", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			if r.URL.Path != "/v1/sessions/mb_destroy1" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		defer srv.Close()

		err := client.DestroySession(context.Background(), "mb_destroy1", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("with save profile", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["save_profile"] != true {
				t.Errorf("expected save_profile=true")
			}
			w.WriteHeader(http.StatusNoContent)
		})
		defer srv.Close()

		err := client.DestroySession(context.Background(), "mb_destroy2", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestAPIError(t *testing.T) {
	t.Run("401", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error":{"code":"unauthorized","message":"Invalid API key"}}`)
		})
		defer srv.Close()

		_, err := client.GetSession(context.Background(), "mb_1")
		if err == nil {
			t.Fatal("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			// Wrapped error
			t.Logf("error: %v", err)
		} else {
			if apiErr.StatusCode != 401 {
				t.Errorf("expected 401, got %d", apiErr.StatusCode)
			}
		}
	})

	t.Run("404", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error":{"code":"not_found","message":"session not found"}}`)
		})
		defer srv.Close()

		_, err := client.GetSession(context.Background(), "nonexistent")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
