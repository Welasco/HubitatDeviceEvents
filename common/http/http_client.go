package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
)

func Http_client(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error("[http][Http_client] Fail to create HTTP request")
		logger.Error("[http][Http_client] Error: " + err.Error())
		return "", err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("[http][Http_client] Failed to connect to: " + url)
		logger.Error("[http][Http_client] Error: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Error("[http][Http_client] Invalid status code from " + url + ": Status Code: " + resp.Status)
		return "", errors.New("Invalid status code from " + url + ": Status Code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("[http][Http_client] Error reading response: " + err.Error())
		return "", err
	}
	bodystr := string(body)
	logger.Debug("[http][Http_client] HTTP Boddy of URL", url, ":\n", bodystr)

	return bodystr, err
}
