package xhttp

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
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
	verifyAPIKey func(ctx context.Context, header http.Header) bool,
	authenticate func(ctx context.Context, header http.Header) *xtype.AuthResult,
	nextHandlers ...http.Handler) http.Handler {
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
				debug.PrintStack()
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

		if verifyAPIKey != nil {
			verified := verifyAPIKey(ctx, req.Header)
			if !verified {
				Error(w, xerror.NewAPIError(http.StatusBadRequest, "invalid api key"))
				return
			}
		}

		if authenticate != nil {
			auth := authenticate(ctx, req.Header)
			if auth != nil {
				if auth.AppID != appID {
					logger.WarnContext(ctx, "invalid app id", slog.String("appId", appID))
				} else {
					a.SetUserID(auth.UserID)
					logger.InfoContext(ctx, "authenticated",
						slog.Any("uid", auth.UserID.Value()),
						slog.String("appId", appID))
				}
			}
		}

		for _, f := range nextHandlers {
			f.ServeHTTP(w, req)
			if w.Status() != 0 {
				return
			}
		}
	})
}

func NewConsumerHandler[T any](f func(ctx context.Context, t T) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params T
		if ReadJSONBody(w, r, &params) {
			err := f(r.Context(), params)
			Error(w, err)
		}
	})
}

func NewSupplierHandler[T any](f func(ctx context.Context) (T, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r.Context())
		JSONOrError(w, res, err)
	})
}

func NewFunctionHandler[T, R any](f func(ctx context.Context, t T) (R, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params T
		if ReadJSONBody(w, r, &params) {
			res, err := f(r.Context(), params)
			JSONOrError(w, res, err)
		}
	})
}
