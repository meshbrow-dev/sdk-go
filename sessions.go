package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
)

// Session represents a browser session.
type Session struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	CDPEndpoint string `json:"cdp_endpoint,omitempty"`
	Token       string `json:"token,omitempty"`
	Stealth     string `json:"stealth,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

// CreateSessionParams configures a new session.
type CreateSessionParams struct {
	Stealth      string    `json:"stealth,omitempty"`
	ProxyType    string    `json:"-"`
	ProxyCountry string    `json:"-"`
	ProfileID    string    `json:"profile_id,omitempty"`
	Viewport     *Viewport `json:"viewport,omitempty"`
}

// Viewport defines browser viewport dimensions.
type Viewport struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// CreateSession launches a new stealth browser session.
func (c *Client) CreateSession(ctx context.Context, params *CreateSessionParams) (*Session, error) {
	body := map[string]any{"stealth": "max"}
	if params != nil {
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
		if params.ProfileID != "" {
			body["profile_id"] = params.ProfileID
		}
		if params.Viewport != nil {
			body["viewport"] = params.Viewport
		}
	}

	data, err := c.do(ctx, "POST", "/v1/sessions", body)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("parse session: %w", err)
	}
	return &sess, nil
}

// GetSession returns details about a session.
func (c *Client) GetSession(ctx context.Context, id string) (*Session, error) {
	data, err := c.do(ctx, "GET", "/v1/sessions/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("parse session: %w", err)
	}
	return &sess, nil
}

// SessionListResponse contains sessions and usage metrics.
type SessionListResponse struct {
	Sessions []Session      `json:"sessions"`
	Metrics  map[string]any `json:"metrics"`
}

// ListSessions returns all active sessions.
func (c *Client) ListSessions(ctx context.Context) (*SessionListResponse, error) {
	data, err := c.do(ctx, "GET", "/v1/sessions", nil)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	var resp SessionListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse sessions: %w", err)
	}
	return &resp, nil
}

// DestroySession destroys a browser session.
func (c *Client) DestroySession(ctx context.Context, id string, saveProfile bool) error {
	var body any
	if saveProfile {
		body = map[string]any{"save_profile": true}
	}
	_, err := c.do(ctx, "DELETE", "/v1/sessions/"+id, body)
	return err
}
