package ai

import (
	"fmt"
	"strings"
)

var blockedTerms = []string{
	"下注", "赔率", "盘口", "博彩", "稳赚", "稳赢", "必中", "100%晋级", "100% 晋级",
	"非法直播", "盗播",
	"odds", "betting", "wager", "bookmaker", "sure win", "lock pick",
}

func ValidateOutput(content string) error {
	lower := strings.ToLower(content)
	for _, term := range blockedTerms {
		if strings.Contains(lower, strings.ToLower(term)) {
			return fmt.Errorf("AI output contains unsupported content")
		}
	}
	return nil
}

func SanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if len(msg) > 240 {
		msg = msg[:240]
	}
	return strings.ReplaceAll(msg, "\n", " ")
}
