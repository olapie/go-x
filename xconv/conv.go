package xconv

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
)

func ToEmailAddress(s string) (string, error) {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func ToURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if u.Scheme == "" {
		return "", errors.New("missing schema")
	}
	if u.Host == "" {
		return "", errors.New("missing host")
	}
	return u.String(), nil
}

const (
	kilobyte = 1 << (10 * (1 + iota)) // 1 << (10*1)
	megabyte                          // 1 << (10*2)
	gigabyte                          // 1 << (10*3)
	terabyte                          // 1 << (10*4)
	petabyte                          // 1 << (10*5)
)

func SizeToHumanReadable(size int64) string {
	if size < kilobyte {
		return fmt.Sprintf("%d B", size)
	} else if size < megabyte {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(kilobyte))
	} else if size < gigabyte {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(megabyte))
	} else if size < terabyte {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(gigabyte))
	} else if size < petabyte {
		return fmt.Sprintf("%.2f TB", float64(size)/float64(terabyte))
	} else {
		return fmt.Sprintf("%.2f PB", float64(size)/float64(petabyte))
	}
}

func Pointer[T any](v T) *T {
	return &v
}
