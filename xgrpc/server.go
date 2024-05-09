package xgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"go.olapie.com/security/base62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xtype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PreHandleOptions struct {
	APIKeyVerifierFunc func(ctx context.Context, md metadata.MD) bool
	AuthenticatorFunc  func(ctx context.Context, md metadata.MD) *xtype.AuthResult
	HeadersRequired    []string
	HeadersForLogging  []string
}

var defaultPreHandleOptions = PreHandleOptions{
	HeadersForLogging: []string{xhttpheader.LowerKeyClientID, xhttpheader.LowerKeyAppID},
}

func PreHandle(ctx context.Context, info *grpc.UnaryServerInfo, optFns ...func(options *PreHandleOptions)) (context.Context, error) {
	options := defaultPreHandleOptions
	for _, fn := range optFns {
		fn(&options)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "failed to read request metadata")
	}

	for _, h := range options.HeadersRequired {
		if getMetadataValue(md, h) == "" {
			return ctx, status.Error(codes.InvalidArgument, "missing "+h)
		}
	}
	traceID := getMetadataValue(md, xhttpheader.LowerKeyTraceID)
	if traceID == "" {
		traceID = base62.NewUUIDString()
	}
	activity := xcontext.NewActivity(info.FullMethod, md)
	activity.SetTraceID(traceID)
	logger := xlog.FromContext(ctx).With(slog.String("trace_id", traceID))
	ctx = xlog.NewContext(xcontext.WithIncomingActivity(ctx, activity), logger)
	logger = xlog.FromContext(ctx).With("module", "xgrpc")
	var logAttrs []any
	for _, h := range options.HeadersForLogging {
		if v := getMetadataValue(md, h); v != "" {
			logAttrs = append(logAttrs, slog.String(h, v))
		}
	}
	logger.InfoContext(ctx, "START", logAttrs...)

	if options.APIKeyVerifierFunc != nil && !options.APIKeyVerifierFunc(ctx, md) {
		logger.ErrorContext(ctx, "invalid api key", slog.Any("metadata", md))
		return ctx, status.Error(codes.InvalidArgument, "failed to verify api key")
	}

	if options.AuthenticatorFunc != nil {
		auth := options.AuthenticatorFunc(ctx, md)
		if auth != nil {
			appID := activity.GetAppID()
			if auth.AppID != appID {
				logger.ErrorContext(ctx, fmt.Sprintf("client app_id %s does not match authenticated app_id %s", appID, auth.AppID))
				return ctx, status.Error(codes.Unauthenticated, "client app_id does not match authenticated app_id")
			}
			activity.SetUserID(auth.UserID)
			logger.InfoContext(ctx, "authenticated", slog.Any("uid", auth.UserID.Value()), slog.String("app_id", auth.AppID))
		}
	}

	return ctx, nil
}

func PostHandle(ctx context.Context, resp any, err error, logger *slog.Logger, startAt time.Time) (any, error) {
	if err == nil {
		logger.InfoContext(ctx, "END", slog.Duration("cost", time.Now().Sub(startAt)))
		return resp, nil
	}

	if reflect.TypeOf(err) == statusErrorType {
		logger.ErrorContext(ctx, "END", xlog.Err(err))
		return nil, err
	}

	if s, ok := status.FromError(err); ok {
		logger.ErrorContext(ctx, "END", slog.Any("status", s), xlog.Err(err))
		return nil, err
	}

	if s := xerror.GetCode(err); s >= 100 && s < 600 {
		code := HTTPStatusToCode(s)
		logger.InfoContext(ctx, "END", slog.Int("status", s), slog.Int("code", int(code)), xlog.Err(err))
		return nil, status.Error(code, err.Error())
	}

	if errors.Is(err, xerror.DBNoRecords) {
		err = status.Error(codes.NotFound, err.Error())
	}
	logger.ErrorContext(ctx, "END", xlog.Err(err))
	return nil, err
}

func getMetadataValue(md metadata.MD, k string) string {
	a := md.Get(k)
	if len(a) == 0 {
		return ""
	}
	return a[0]
}
