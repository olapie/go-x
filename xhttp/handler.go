package xhttp

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.olapie.com/x/xbase62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xtype"
)

type joinHandler struct {
	handlers []http.Handler
}

func (j *joinHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	w := WrapWriter(writer)
	for _, h := range j.handlers {
		h.ServeHTTP(w, request)
		if w.Status() != 0 {
			return
		}
	}
}

var _ http.Handler = (*joinHandler)(nil)

func JoinHandlers(handlers ...http.Handler) http.Handler {
	return &joinHandler{
		handlers: handlers,
	}
}

func JoinHandlerFuncs(funcs ...http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		w := WrapWriter(writer)
		for _, f := range funcs {
			f.ServeHTTP(w, request)
			if w.Status() != 0 {
				return
			}
		}
	}
}

type Authenticator[T xtype.UserIDTypes] interface {
	Authenticate(req *http.Request) (T, error)
}

func NewStartHandler(
	maybeNext http.Handler,
	verifyAPIKey func(ctx context.Context, header http.Header) bool,
	authenticate func(ctx context.Context, header http.Header) *xtype.AuthResult) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		startAt := time.Now()
		a := xcontext.NewActivity("", req.Header)
		ctx := xcontext.WithIncomingActivity(req.Context(), a)
		traceID := a.Get(xhttpheader.KeyTraceID)
		if traceID == "" {
			traceID = xbase62.NewUUIDString()
		}

		logger := xlog.FromContext(ctx).With(slog.String("traceId", traceID))

		fields := make([]any, 0, 4+len(req.Header))
		fields = append(fields,
			slog.String("uri", req.RequestURI),
			slog.String("method", req.Method),
			slog.String("host", req.Host),
			slog.String("remoteAddr", req.RemoteAddr))
		for key := range req.Header {
			fields = append(fields, slog.String(key, req.Header.Get(key)))
		}

		ctx = xlog.NewContext(ctx, logger)
		logger = logger.With("module", "xhttp")
		logger.InfoContext(ctx, "START", fields...)
		var cancel context.CancelFunc
		if seconds := a.GetRequestTimeout(); seconds > 0 {
			ctx, cancel = context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
			defer cancel()
		}

		req = req.WithContext(ctx)
		w := WrapWriter(rw)

		defer func() {
			if p := recover(); p != nil {
				logger.ErrorContext(ctx, "panic", "error", p)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			status := w.Status()
			fields = []any{slog.Int("status", status),
				slog.Duration("cost", time.Now().Sub(startAt))}
			if status >= 400 {
				fields = append(fields, slog.String("body", string(w.Body())))
				logger.ErrorContext(ctx, "END", fields...)
			} else {
				logger.InfoContext(ctx, "END", fields...)
			}
		}()

		appID := a.GetAppID()
		if appID == "" {
			Error(w, xerror.NewAPIError(http.StatusBadRequest, "client appId does not match authenticated appId"))
			return
		}

		if verifyAPIKey(ctx, req.Header) {
			auth := authenticate(ctx, req.Header)
			if auth != nil {
				if auth.AppID != appID {
					Error(w, xerror.NewAPIError(http.StatusUnauthorized, "client appId does not match authenticated appId"))
					return
				} else {
					a.SetUserID(auth.UserID)
					logger.InfoContext(ctx, "authenticated", slog.Any("uid", auth.UserID.Value()), slog.String("appId", auth.AppID))
				}
			}
			maybeNext.ServeHTTP(w, req)
		} else {
			Error(w, xerror.NewAPIError(http.StatusBadRequest, "invalid api key"))
			return
		}
	})
}
