package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Device mirrors model.Device from the main application.
// Note: the JSON tag for the identifier is uppercase "Id" to match the Hubitat API.
type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Room  string `json:"room"`
}

// HubitatClient is a thin HTTP client for the HubitatDeviceEvents REST API.
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

// do performs an HTTP request and returns the raw response body. A non-2xx
// status code is returned as an error that includes the status and body so the
// AI client can surface a useful message.
func (c *HubitatClient) do(method, path string, body any) (string, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("encoding request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	return c.do(http.MethodGet, "/api/v1/device", nil)
}

func (c *HubitatClient) GetDevice(id string) (string, error) {
	return c.do(http.MethodGet, "/api/v1/device/"+url.PathEscape(id), nil)
}

func (c *HubitatClient) AddDevice(d Device) (string, error) {
	return c.do(http.MethodPost, "/api/v1/device", d)
}

func (c *HubitatClient) UpdateDevice(d Device) (string, error) {
	return c.do(http.MethodPut, "/api/v1/device", d)
}

func (c *HubitatClient) DeleteDevice(id string) (string, error) {
	return c.do(http.MethodDelete, "/api/v1/device/"+url.PathEscape(id), nil)
}

func (c *HubitatClient) ListDeviceEvents() (string, error) {
	return c.do(http.MethodGet, "/api/v1/device/event", nil)
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
	return c.do(http.MethodGet, path, nil)
}
