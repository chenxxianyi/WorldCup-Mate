package ai

import (
	"encoding/json"
	"strings"
)

func DecodeJSON[T any](raw string, fallback T) (T, string, error) {
	candidate := ExtractJSON(raw)
	if strings.TrimSpace(candidate) == "" {
		return fallback, raw, json.Unmarshal([]byte(raw), &fallback)
	}
	var out T
	if err := json.Unmarshal([]byte(candidate), &out); err != nil {
		return fallback, raw, err
	}
	return out, candidate, nil
}

func ExtractJSON(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}

	objStart := strings.Index(s, "{")
	objEnd := strings.LastIndex(s, "}")
	if objStart >= 0 && objEnd > objStart {
		return s[objStart : objEnd+1]
	}
	arrStart := strings.Index(s, "[")
	arrEnd := strings.LastIndex(s, "]")
	if arrStart >= 0 && arrEnd > arrStart {
		return s[arrStart : arrEnd+1]
	}
	return s
}

func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
