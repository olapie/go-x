package xgrpc

import (
	"context"
	"fmt"
	"go.olapie.com/x/xbase62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xlog"
	"log/slog"
	"reflect"
	"time"

	"go.olapie.com/ola/errorutil"
	"go.olapie.com/ola/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var metadataKeysForLogging = []string{"x-client-id", "x-app-id"}

func ServerStart(ctx context.Context,
	info *grpc.UnaryServerInfo,
	verifyAPIKey func(ctx context.Context, md metadata.MD) bool,
	authenticate func(ctx context.Context, md metadata.MD) *types.Auth) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "failed reading request metadata")
	}

	a := xcontext.NewActivity(info.FullMethod, md)
	appID := a.GetAppID()
	if appID == "" {
		return ctx, status.Error(codes.InvalidArgument, "missing x-app-id")
	}

	ctx = xcontext.WithIncomingActivity(ctx, a)
	traceID := a.GetTraceID()
	if traceID == "" {
		traceID = xbase62.NewUUIDString()
		a.SetTraceID(traceID)
	}
	logger := xlog.FromContext(ctx).With(slog.String("traceId", traceID))
	ctx = xlog.NewContext(ctx, logger)
	logger = logger.With("module", "xgrpc")
	fields := make([]any, 0, len(md)+1)
	fields = append(fields, slog.String("method", info.FullMethod))

	for _, mdKey := range metadataKeysForLogging {
		if mdVal, _ := md[mdKey]; len(mdVal) > 0 && mdVal[0] != "" {
			fields = append(fields, slog.String(mdKey, mdVal[0]))
		}
	}
	logger.Info("START", fields...)

	if !verifyAPIKey(ctx, md) {
		attrs := make([]slog.Attr, 0, len(md))
		for key, val := range md {
			if len(val) > 0 {
				attrs = append(attrs, slog.String(key, val[0]))
			}
		}
		logger.Error("invalid api key", slog.Any("metadata", md))
		return ctx, status.Error(codes.InvalidArgument, "failed verifying")
	}

	auth := authenticate(ctx, md)
	if auth != nil {
		if auth.AppID != appID {
			logger.ErrorContext(ctx, fmt.Sprintf("client appId %s does not match authenticated appId %s", appID, auth.AppID))
			return ctx, status.Error(codes.Unauthenticated, "client appId does not match authenticated appId")
		}
		a.SetUserID(auth.UserID)
		logger.Info("authenticated", slog.Any("uid", auth.UserID.Value()), slog.String("appId", auth.AppID))
	}
	return ctx, nil
}

func ServerFinish(resp any, err error, logger *slog.Logger, startAt time.Time) (any, error) {
	if err == nil {
		logger.Info("END", slog.Duration("cost", time.Now().Sub(startAt)))
		return resp, nil
	}

	if reflect.TypeOf(err) == statusErrorType {
		logger.Error("END", xlog.Err(err))
		return nil, err
	}

	if s := errorutil.GetCode(err); s >= 100 && s < 600 {
		code := HTTPStatusToCode(s)
		logger.Info("END", slog.Int("status", s), slog.Int("code", int(code)), xlog.Err(err))
		return nil, status.Error(code, err.Error())
	}
	logger.Error("END", xlog.Err(err))

	return nil, err
}
