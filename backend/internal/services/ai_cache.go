package services

import "strings"

func sanitizeCachePart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "none"
	}
	replacer := strings.NewReplacer(":", "_", " ", "_", "/", "_", "\\", "_")
	return replacer.Replace(value)
}
