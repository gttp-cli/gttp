package utils

import (
	"io"
	"net/http"
	"os"
)

// ReadURL sends a GET request to the specified URL and returns the response body as a string.
func ReadURL(url string) (string, error) {
	if url[0:4] != "http" {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ReadFile reads the specified file and returns the contents as a string.
func ReadFile(file string) (string, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
