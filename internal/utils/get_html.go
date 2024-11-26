package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error fetching %s: %w", rawURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("error reading URL %s", rawURL)
	}

	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("not HTML content found at %s", rawURL)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading content from %s: %w", rawURL, err)
	}

	return string(data), nil
}
