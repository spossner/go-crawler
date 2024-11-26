package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	normalizedURL, _ := strings.CutSuffix(fmt.Sprintf("%s%s", u.Host, u.Path), "/")
	return normalizedURL, nil
}
