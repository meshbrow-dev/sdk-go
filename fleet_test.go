package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestCreateFleet(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/fleet" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["count"] != float64(5) {
			t.Errorf("expected count=5, got %v", body["count"])
		}
		proxy, ok := body["proxy"].(map[string]any)
		if !ok {
			t.Fatal("expected proxy in body")
		}
		if proxy["type"] != "residential" {
			t.Errorf("expected proxy type residential")
		}
		if proxy["country"] != "DE" {
			t.Errorf("expected proxy country DE")
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"data":{"id":"fleet_1","status":"ready","sessions":[{"id":"mb_1","status":"ready"},{"id":"mb_2","status":"ready"}],"count":2}}`)
	})
	defer srv.Close()

	fleet, err := client.CreateFleet(context.Background(), &CreateFleetParams{
		Count:        5,
		Stealth:      "max",
		ProxyType:    "residential",
		ProxyCountry: "DE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fleet.ID != "fleet_1" {
		t.Errorf("expected ID fleet_1, got %s", fleet.ID)
	}
	if len(fleet.Sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(fleet.Sessions))
	}
}

func TestGetFleet(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/fleet/fleet_abc" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		fmt.Fprintf(w, `{"data":{"id":"fleet_abc","status":"ready","sessions":[{"id":"mb_1","status":"ready"}],"count":1}}`)
	})
	defer srv.Close()

	fleet, err := client.GetFleet(context.Background(), "fleet_abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fleet.ID != "fleet_abc" {
		t.Errorf("expected ID fleet_abc, got %s", fleet.ID)
	}
	if fleet.Count != 1 {
		t.Errorf("expected count 1, got %d", fleet.Count)
	}
}

func TestDestroyFleet(t *testing.T) {
	client, srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/fleet/fleet_abc" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer srv.Close()

	err := client.DestroyFleet(context.Background(), "fleet_abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
