package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestNavigate(t *testing.T) {
	t.Run("basic navigate", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v1/sessions/mb_1/navigate" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["url"] != "https://example.com" {
				t.Errorf("expected url https://example.com, got %v", body["url"])
			}
			if body["wait_until"] != "load" {
				t.Errorf("expected wait_until=load, got %v", body["wait_until"])
			}

			fmt.Fprintf(w, `{"data":{"url":"https://example.com","title":"Example Domain","status":200}}`)
		})
		defer srv.Close()

		resp, err := client.Navigate(context.Background(), "mb_1", "https://example.com", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.URL != "https://example.com" {
			t.Errorf("expected url https://example.com, got %s", resp.URL)
		}
		if resp.Title != "Example Domain" {
			t.Errorf("expected title Example Domain, got %s", resp.Title)
		}
		if resp.Status != 200 {
			t.Errorf("expected status 200, got %d", resp.Status)
		}
	})

	t.Run("with networkidle", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["wait_until"] != "networkidle" {
				t.Errorf("expected wait_until=networkidle, got %v", body["wait_until"])
			}
			fmt.Fprintf(w, `{"data":{"url":"https://example.com","title":"Example","status":200}}`)
		})
		defer srv.Close()

		_, err := client.Navigate(context.Background(), "mb_1", "https://example.com", "networkidle")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestTakeScreenshot(t *testing.T) {
	t.Run("full page", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sessions/mb_1/screenshot" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["full_page"] != true {
				t.Errorf("expected full_page=true")
			}

			fmt.Fprintf(w, `{"data":{"data":"iVBORw0KGgo=","format":"png"}}`)
		})
		defer srv.Close()

		s, err := client.TakeScreenshot(context.Background(), "mb_1", "", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Data != "iVBORw0KGgo=" {
			t.Errorf("unexpected data: %s", s.Data)
		}
		if s.Format != "png" {
			t.Errorf("expected format png, got %s", s.Format)
		}
	})

	t.Run("with selector", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["selector"] != "#main" {
				t.Errorf("expected selector=#main, got %v", body["selector"])
			}
			fmt.Fprintf(w, `{"data":{"data":"abc123","format":"png"}}`)
		})
		defer srv.Close()

		_, err := client.TakeScreenshot(context.Background(), "mb_1", "#main", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestClick(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sessions/mb_1/click" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["selector"] != "button.submit" {
			t.Errorf("expected selector=button.submit, got %v", body["selector"])
		}
		fmt.Fprintf(w, `{"data":{"clicked":true}}`)
	})
	defer srv.Close()

	err := client.Click(context.Background(), "mb_1", "button.submit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestType(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["selector"] != "input#email" {
			t.Errorf("expected selector=input#email, got %v", body["selector"])
		}
		if body["text"] != "test@example.com" {
			t.Errorf("unexpected text: %v", body["text"])
		}
		if body["clear"] != true {
			t.Errorf("expected clear=true")
		}
		fmt.Fprintf(w, `{"data":{"typed":true}}`)
	})
	defer srv.Close()

	err := client.Type(context.Background(), "mb_1", "input#email", "test@example.com", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtract(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"data":{"text":"Hello World"}}`)
		})
		defer srv.Close()

		resp, err := client.Extract(context.Background(), "mb_1", "", 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Text != "Hello World" {
			t.Errorf("expected Hello World, got %s", resp.Text)
		}
	})

	t.Run("with selector and max_length", func(t *testing.T) {
		client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["selector"] != "article" {
				t.Errorf("expected selector=article, got %v", body["selector"])
			}
			if body["max_length"] != float64(1000) {
				t.Errorf("expected max_length=1000, got %v", body["max_length"])
			}
			fmt.Fprintf(w, `{"data":{"text":"Article content"}}`)
		})
		defer srv.Close()

		resp, err := client.Extract(context.Background(), "mb_1", "article", 1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Text != "Article content" {
			t.Errorf("unexpected text: %s", resp.Text)
		}
	})
}

func TestExecute(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["script"] != "return 21 * 2" {
			t.Errorf("unexpected script: %v", body["script"])
		}
		fmt.Fprintf(w, `{"data":{"result":42}}`)
	})
	defer srv.Close()

	resp, err := client.Execute(context.Background(), "mb_1", "return 21 * 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Result != float64(42) {
		t.Errorf("expected 42, got %v", resp.Result)
	}
}
