package xhttpheader

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"go.olapie.com/x/xerror"
)

var getCurrentTime = func() time.Time {
	return time.Now()
}

func SetCurrentTimeFunc(f func() time.Time) {
	if f != nil {
		getCurrentTime = f
	}
}

func Sign[T ~map[string][]string](h T) {
	clientID := Get(h, KeyClientID)
	if clientID == "" {
		panic("missing " + KeyClientID)
	}
	traceID := Get(h, KeyTraceID)
	if traceID == "" {
		panic("missing " + KeyTraceID)
	}
	timestamp := Get(h, KeyTimestamp)
	if timestamp == "" {
		timestamp = fmt.Sprint(getCurrentTime().Unix())
		Set(h, KeyTimestamp, timestamp)
	}
	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic("invalid " + timestamp)
	}
	sign := generateAPIKey(clientID, traceID, t)
	Set(h, KeyAPIKey, sign)
}

func Verify[T ~map[string][]string | ~map[string]string](h T, leeway time.Duration) error {
	sign := Get(h, KeyAPIKey)
	if sign == "" {
		return xerror.BadRequest("missing " + KeyAPIKey)
	}

	clientID := Get(h, KeyClientID)
	if clientID == "" {
		return xerror.BadRequest("missing " + KeyClientID)
	}
	traceID := Get(h, KeyTraceID)
	if traceID == "" {
		return xerror.BadRequest("missing " + KeyTraceID)
	}
	timestamp := Get(h, KeyTimestamp)
	if timestamp == "" {
		return xerror.BadRequest("missing " + KeyTimestamp)
	}

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return xerror.BadRequest("invalid %s: %s", KeyTimestamp, timestamp)
	}

	if time.Unix(t, 0).Add(leeway).Before(getCurrentTime()) {
		return xerror.BadRequest("expired %s: %s", KeyAPIKey, timestamp)
	}

	expected := generateAPIKey(clientID, traceID, t)
	if sign != expected {
		return xerror.BadRequest("invalid " + KeyAPIKey)
	}
	return nil
}

func generateAPIKey(clientID, traceID string, timestamp int64) string {
	var digest = fmt.Sprintf("%s.%s.%d", clientID, traceID, timestamp)
	return fmt.Sprintf("%x", md5.Sum([]byte(digest)))
}
