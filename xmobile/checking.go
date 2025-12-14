package xmobile

import (
	"net/mail"
	"strings"
	"unicode"

	"go.olapie.com/times"
	"go.olapie.com/x/xurl"
)

func IsEmailAddress(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func IsURL(s string) bool {
	return xurl.IsURL(s)
}

func IsDate(s string) bool {
	return times.IsDate(s)
}

var (
	MinPasswordLen int = 6
	MinUsernameLen int = 4
	MaxUsernameLen int = 20
)

func IsValidPassword(password string) bool {
	if len(password) < MinPasswordLen {
		return false
	}

	hasDigit := false
	hasAlpha := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasAlpha = true
		}
	}
	return hasDigit && hasAlpha
}

func IsValidUsername(username string) bool {
	username = strings.ToLower(username)
	s := []rune(username)
	if len(s) > MaxUsernameLen {
		return false
	}

	if len(s) < MinUsernameLen {
		return false
	}

	for _, c := range s {
		if unicode.IsDigit(c) || c == '_' || (c >= 'a' && c <= 'z') {
			continue
		}
		return false
	}
	return true
}
