package http

import (
	"errors"
	"io"
	"net/http"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
)

func Http_client(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("[http][Http_client] Fail to create HTTP request")
		logger.Error(err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("[http][Http_client] Failed to connect to: " + url)
		logger.Error("[http][Http_client] Error: " + err.Error())
		return "", err
	}
	if resp.StatusCode != 200 {
		logger.Error("[http][Http_client] Invalid status code from " + url + ": Status Code: " + resp.Status)
		return "", errors.New("Invalid status code from " + url + ": Status Code: " + resp.Status)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	logger.Debug("[http][Http_client] HTTP Boddy of URL", url, ":\n", bodystr)

	return bodystr, err
}
