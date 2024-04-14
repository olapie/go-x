package xcontext

import (
	"net/http"
	"testing"
)

func TestNewActivity(t *testing.T) {
	t.Run("NewWithHTTPHeader", func(t *testing.T) {
		NewActivity("http", http.Header{})
	})

	type Map map[string]string

	t.Run("NewWithMap", func(t *testing.T) {
		NewActivity("http", Map{})
	})
}
