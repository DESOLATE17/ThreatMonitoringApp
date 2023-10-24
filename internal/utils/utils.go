package utils

import "strings"

func ExtractObjectNameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
