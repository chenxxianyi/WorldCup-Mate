package utils

import (
	"errors"
	"regexp"
	"strings"
)

// Password policy (SEC-03): min 8 chars, at least one letter and one digit,
// and not in the common-weak-password list. Shared by register and
// change-password; the server-side check is authoritative.
var (
	ErrWeakPassword = errors.New("密码至少 8 位，且需同时包含字母和数字")
)

var (
	hasLetterRe = regexp.MustCompile(`[A-Za-z]`)
	hasDigitRe  = regexp.MustCompile(`[0-9]`)
)

var commonPasswords = map[string]bool{
	"12345678": true, "123456789": true, "1234567890": true,
	"password": true, "password1": true, "passw0rd": true,
	"qwerty123": true, "qwertyui": true, "iloveyou1": true,
	"admin123": true, "admin1234": true, "abc12345": true,
	"1234567a": true, "a1234567": true, "88888888": true,
	"11111111": true, "00000000": true, "welcome1": true,
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	if !hasLetterRe.MatchString(password) || !hasDigitRe.MatchString(password) {
		return ErrWeakPassword
	}
	if commonPasswords[strings.ToLower(password)] {
		return ErrWeakPassword
	}
	return nil
}
