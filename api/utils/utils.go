package utils

import (
	"os"
	"strings"
)

func IsDifferentDomain(url string) bool {

	domain := os.Getenv("DOMAIN")

	if url == domain {
		return false
	}

	if strings.Contains(url, domain) {
		return false
	}

	cleanURL := strings.TrimPrefix(url, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "www.")
	cleanURL = strings.Split(cleanURL, "/")[0]

	return cleanURL != domain

}

func EnsureHttpPrefix(url string) string {
	if !strings.HasPrefix(url, "https") || !strings.HasPrefix(url, "http") {
		return "https://" + url
	}
	return url
}
