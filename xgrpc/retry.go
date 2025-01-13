package xgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type CallFunc[IN proto.Message, OUT proto.Message] func(ctx context.Context, in IN, options ...grpc.CallOption) (OUT, error)

type RetryOptions struct {
	Count              int
	Backoff            time.Duration
	RefreshAccessToken func(ctx context.Context) (string, error)
}

type Retry[IN proto.Message, OUT proto.Message] struct {
	options  RetryOptions
	call     CallFunc[IN, OUT]
	callName string
}

func NewRetry[IN proto.Message, OUT proto.Message](call CallFunc[IN, OUT], options ...func(options *RetryOptions)) *Retry[IN, OUT] {
	r := &Retry[IN, OUT]{
		call:     call,
		callName: xreflect.FuncNameOf(call),
	}
	for _, opt := range options {
		opt(&r.options)
	}
	if r.options.Count <= 0 {
		r.options.Count = 3
	}
	if r.options.Backoff <= 0 {
		r.options.Backoff = time.Second // grpc conn initial backoff is one second, use the same value here
	}
	return r
}

func (r *Retry[IN, OUT]) Call(ctx context.Context, in IN, options ...grpc.CallOption) (OUT, error) {
	var out OUT
	var err error
	logger := xlog.FromContext(ctx).With("module", "xgrpc").With("call", r.callName)
	for i := 0; i < r.options.Count; i++ {
		if i > 0 {
			logger.Debug(fmt.Sprintf("retrying %d", i))
		}
		out, err = r.call(ctx, in, options...)
		if err == nil {
			return out, nil
		}

		if errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
			return out, fmt.Errorf("failed to call due to context error: %w", err)
		}

		switch code := GetErrorCode(err); code {
		// unrecoverable error codes, return immediately
		case codes.InvalidArgument,
			codes.Unimplemented,
			codes.PermissionDenied,
			codes.Internal,
			codes.AlreadyExists,
			codes.NotFound,
			codes.DeadlineExceeded:
			return out, fmt.Errorf("failed to call: %w", err)
		case codes.Unauthenticated:
			if r.options.RefreshAccessToken == nil {
				return out, fmt.Errorf("failed to call: %w, and options.RefreshAccessToken is nil", err)
			}
			logger.DebugContext(ctx, "refreshing access token")
			accessToken, err := r.options.RefreshAccessToken(ctx)

			if err == nil {
				act := xcontext.GetOutgoingActivity(ctx)
				if act == nil {
					return out, xerror.BadRequest("no outgoing activity")
				}
				act.SetAuthorization(accessToken)
				logger.DebugContext(ctx, "refreshed access token successfully")
			} else {
				if errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
					return out, fmt.Errorf("failed to refresh access token due to context error: %w", err)
				}

				switch refreshAccessTokenErrorCode := GetErrorCode(err); refreshAccessTokenErrorCode {
				case codes.InvalidArgument,
					codes.Unimplemented,
					codes.PermissionDenied,
					codes.Internal,
					codes.AlreadyExists,
					codes.NotFound,
					codes.DeadlineExceeded,
					codes.Unauthenticated:
					return out, fmt.Errorf("failed to refresh access token: %w", err)
				default:
					logger.Error("failed to refresh access token", slog.Int("code", int(refreshAccessTokenErrorCode)), xlog.Err(err))
				}
			}
		default:
			time.Sleep(r.options.Backoff)
		}
	}
	return out, err
}
