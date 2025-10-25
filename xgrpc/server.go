package xgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"time"

	"go.olapie.com/x/xbase62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xtype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ServerInterceptorOptions struct {
	APIKeyVerifierFunc    func(ctx context.Context, md metadata.MD, fullMethod string) bool
	AuthenticatorFunc     func(ctx context.Context, md metadata.MD, fullMethod string) *xtype.AuthResult
	MetadataValidatorFunc func(ctx context.Context, md metadata.MD, fullMethod string) error
	RequiredMetadataKeys  []string
	LoggingMetadataKeys   []string
}

func NewServerInterceptor(optFns ...func(options *ServerInterceptorOptions)) grpc.UnaryServerInterceptor {
	options := &ServerInterceptorOptions{
		LoggingMetadataKeys: []string{xhttpheader.LowerKeyClientID, xhttpheader.LowerKeyAppID},
	}
	for _, fn := range optFns {
		fn(options)
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		startAt := time.Now()
		logger := xlog.FromContext(ctx).With("module", "xgrpc")
		defer func() {
			if e := recover(); e != nil {
				logger.Error("recovered from a panic", slog.Any("panic", e), slog.String("stack", string(debug.Stack())))
				err = xerror.InternalServerError("")
			}
		}()
		ctx, err = preprocess(ctx, info, options)
		logger = xlog.FromContext(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "failed to preprocess request", xlog.Err(err))
			return nil, err
		}

		if logger.Enabled(ctx, slog.LevelDebug) {
			if msg, ok := req.(proto.Message); ok {
				logger.Debug("processing request", slog.Any("request", protojson.MarshalOptions{}.Format(msg)))
			} else {
				logger.Debug("processing request", slog.Any("request", req))
			}
		}

		resp, err = handler(ctx, req)
		return postprocess(ctx, resp, err, logger, startAt)
	}
}

func GetMetadataValue(md metadata.MD, k string) string {
	a := md.Get(k)
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func preprocess(ctx context.Context, info *grpc.UnaryServerInfo, options *ServerInterceptorOptions) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "failed to read request metadata")
	}

	traceID := GetMetadataValue(md, xhttpheader.LowerKeyTraceID)
	if traceID == "" {
		traceID = xbase62.NewUUIDString()
	}
	activity := xcontext.NewActivity(info.FullMethod, md)
	activity.SetTraceID(traceID)
	logger := xlog.FromContext(ctx).With(slog.String("trace_id", traceID))
	ctx = xlog.NewContext(xcontext.WithIncomingActivity(ctx, activity), logger)
	var logAttrs []any
	for _, h := range options.LoggingMetadataKeys {
		if v := GetMetadataValue(md, h); v != "" {
			logAttrs = append(logAttrs, slog.String(h, v))
		}
	}
	logger.InfoContext(ctx, "START", logAttrs...)

	for _, h := range options.RequiredMetadataKeys {
		if GetMetadataValue(md, h) == "" {
			return ctx, status.Error(codes.InvalidArgument, "missing "+h)
		}
	}

	if options.MetadataValidatorFunc != nil {
		if err := options.MetadataValidatorFunc(ctx, md, info.FullMethod); err != nil {
			return ctx, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to validate metadata: %v", err))
		}
	}

	if options.APIKeyVerifierFunc != nil && !options.APIKeyVerifierFunc(ctx, md, info.FullMethod) {
		return ctx, status.Error(codes.InvalidArgument, "missing or invalid api key")
	}

	if options.AuthenticatorFunc != nil {
		auth := options.AuthenticatorFunc(ctx, md, info.FullMethod)
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

func postprocess(ctx context.Context, resp any, err error, logger *slog.Logger, startAt time.Time) (any, error) {
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
