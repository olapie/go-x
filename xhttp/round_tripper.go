package xhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt == nil {
		return http.DefaultTransport.RoundTrip(r)
	}
	return rt(r)
}

func LogRoundTripper(logger *slog.Logger, next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		logger.Info(fmt.Sprintf("CLIENT Request: %s %s %v", r.Method, r.URL, r.Header))
		start := time.Now()

		resp, err := next.RoundTrip(r)

		duration := time.Since(start)
		if err != nil {
			logger.Error(fmt.Sprintf("CLIENT Request failed after %s: %v", duration, err))
			return nil, err
		}
		logger.Info(fmt.Sprintf("CLIENT Response: %s %s %d in %s", r.Method, r.URL, resp.StatusCode, duration))
		return resp, nil
	})
}

func AuthorizationRoundTripper(refreshAuthorizationFunc func(ctx context.Context) (string, error), next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			auth, refreshErr := refreshAuthorizationFunc(req.Context())
			if refreshErr != nil {
				return nil, err
			}

			req.Header.Set(xhttpheader.KeyAuthorization, auth)
			return next.RoundTrip(req)
		}

		return resp, nil
	})
}

func StatusCheckRoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 400 {
			var buf bytes.Buffer
			_, ioErr := io.Copy(&buf, resp.Body)
			if ioErr != nil {
				return resp, ioErr
			}
			resp.Body = io.NopCloser(&buf)
			return resp, xerror.NewAPIError(resp.StatusCode, buf.String())
		}
		return resp, err
	})
}
