package meshbrow

import (
	"context"
	"encoding/json"
	"fmt"
)

// NavigateResponse contains the result of a navigation.
type NavigateResponse struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

// Screenshot contains a base64-encoded screenshot.
type Screenshot struct {
	Data   string `json:"data"`   // base64 PNG
	Format string `json:"format"` // "png"
}

// ExtractResponse contains extracted text.
type ExtractResponse struct {
	Text string `json:"text"`
}

// ExecuteResponse contains JS execution result.
type ExecuteResponse struct {
	Result any `json:"result"`
}

// Navigate navigates the browser to a URL.
func (c *Client) Navigate(ctx context.Context, sessionID, url string, waitUntil string) (*NavigateResponse, error) {
	if waitUntil == "" {
		waitUntil = "load"
	}
	data, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/navigate", map[string]any{
		"url":        url,
		"wait_until": waitUntil,
	})
	if err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}

	var resp NavigateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse navigate: %w", err)
	}
	return &resp, nil
}

// TakeScreenshot captures a screenshot.
func (c *Client) TakeScreenshot(ctx context.Context, sessionID string, selector string, fullPage bool) (*Screenshot, error) {
	body := map[string]any{"full_page": fullPage}
	if selector != "" {
		body["selector"] = selector
	}

	data, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/screenshot", body)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}

	var s Screenshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse screenshot: %w", err)
	}
	return &s, nil
}

// Click clicks an element.
func (c *Client) Click(ctx context.Context, sessionID, selector string) error {
	_, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/click", map[string]any{
		"selector": selector,
	})
	return err
}

// Type types text into an input field.
func (c *Client) Type(ctx context.Context, sessionID, selector, text string, clear bool) error {
	_, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/type", map[string]any{
		"selector": selector,
		"text":     text,
		"clear":    clear,
	})
	return err
}

// Extract extracts text content from the page.
func (c *Client) Extract(ctx context.Context, sessionID string, selector string, maxLength int) (*ExtractResponse, error) {
	body := map[string]any{}
	if selector != "" {
		body["selector"] = selector
	}
	if maxLength > 0 {
		body["max_length"] = maxLength
	}

	data, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/extract", body)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	var resp ExtractResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse extract: %w", err)
	}
	return &resp, nil
}

// Execute runs JavaScript in the page context.
func (c *Client) Execute(ctx context.Context, sessionID, script string) (*ExecuteResponse, error) {
	data, err := c.do(ctx, "POST", "/v1/sessions/"+sessionID+"/execute", map[string]any{
		"script": script,
	})
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	var resp ExecuteResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse execute: %w", err)
	}
	return &resp, nil
}
