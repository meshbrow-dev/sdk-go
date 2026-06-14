package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
)

// Fleet represents a group of browser sessions.
type Fleet struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	Sessions []Session `json:"sessions"`
	Count    int       `json:"count"`
}

// CreateFleetParams configures a new fleet.
type CreateFleetParams struct {
	Count        int    `json:"count"`
	Stealth      string `json:"stealth,omitempty"`
	ProxyType    string `json:"-"`
	ProxyCountry string `json:"-"`
}

// CreateFleet launches multiple sessions in parallel.
func (c *Client) CreateFleet(ctx context.Context, params *CreateFleetParams) (*Fleet, error) {
	body := map[string]any{
		"count":   params.Count,
		"stealth": "max",
	}
	if params.Stealth != "" {
		body["stealth"] = params.Stealth
	}
	if params.ProxyType != "" {
		proxy := map[string]any{"type": params.ProxyType}
		if params.ProxyCountry != "" {
			proxy["country"] = params.ProxyCountry
		}
		body["proxy"] = proxy
	}

	data, err := c.do(ctx, "POST", "/v1/fleet", body)
	if err != nil {
		return nil, fmt.Errorf("create fleet: %w", err)
	}

	var fleet Fleet
	if err := json.Unmarshal(data, &fleet); err != nil {
		return nil, fmt.Errorf("parse fleet: %w", err)
	}
	return &fleet, nil
}

// GetFleet returns fleet status.
func (c *Client) GetFleet(ctx context.Context, id string) (*Fleet, error) {
	data, err := c.do(ctx, "GET", "/v1/fleet/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("get fleet: %w", err)
	}

	var fleet Fleet
	if err := json.Unmarshal(data, &fleet); err != nil {
		return nil, fmt.Errorf("parse fleet: %w", err)
	}
	return &fleet, nil
}

// DestroyFleet destroys all sessions in a fleet.
func (c *Client) DestroyFleet(ctx context.Context, id string) error {
	_, err := c.do(ctx, "DELETE", "/v1/fleet/"+id, nil)
	return err
}
