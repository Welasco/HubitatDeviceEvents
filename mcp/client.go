package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HubitatClient is a thin, read-only HTTP client for the HubitatDeviceEvents REST API.
type HubitatClient struct {
	baseURL string
	http    *http.Client
}

func NewHubitatClient(baseURL string) *HubitatClient {
	return &HubitatClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// get performs an HTTP GET and returns the raw response body. A non-2xx status
// code is returned as an error that includes the status and body so the AI
// client can surface a useful message.
func (c *HubitatClient) get(path string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("connecting to %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API returned %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	return string(respBody), nil
}

func (c *HubitatClient) ListDevices() (string, error) {
	return c.get("/api/v1/device")
}

func (c *HubitatClient) GetDevice(id string) (string, error) {
	return c.get("/api/v1/device/" + url.PathEscape(id))
}

func (c *HubitatClient) ListDeviceEvents() (string, error) {
	return c.get("/api/v1/device/event")
}

func (c *HubitatClient) GetDeviceEvents(id, start, end string) (string, error) {
	path := "/api/v1/device/" + url.PathEscape(id) + "/event"
	q := url.Values{}
	if start != "" {
		q.Set("start", start)
	}
	if end != "" {
		q.Set("end", end)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return c.get(path)
}
