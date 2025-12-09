package xgrpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"go.olapie.com/x/xbase62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DialTLS(ctx context.Context, server string, options ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	return grpc.DialContext(ctx, server, options...)
}

func DialWithClientCert(ctx context.Context, server string, cert []byte, options ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	config := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
			},
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	return grpc.DialContext(ctx, server, options...)
}

func Dial(ctx context.Context, server string, options ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.DialContext(ctx, server, options...)
}

// WithSigner set trace id, api key and other properties in metadata
func WithSigner(createAPIKey func(md metadata.MD)) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = signClientContext(ctx, createAPIKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}

func signClientContext(ctx context.Context, createAPIKey func(md metadata.MD)) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}

	a := xcontext.GetOutgoingActivity(ctx)
	if a != nil {
		xcontext.CopyActivityHeader(md, a)
	} else {
		xlog.FromContext(ctx).Warn("no outgoing context")
	}
	if traceID := xhttpheader.GetTraceID(md); traceID == "" {
		if traceID == "" {
			traceID = xbase62.NewUUIDString()
			xlog.FromContext(ctx).InfoContext(ctx, "xgrpc.signClientContext: generated trace id "+traceID)
		}
		xhttpheader.SetTraceID(md, traceID)
	}
	if timestamp := xhttpheader.Get(md, xhttpheader.LowerKeyTimestamp); timestamp == "" {
		timestamp = fmt.Sprint(time.Now().UnixMilli())
		xhttpheader.Set(md, xhttpheader.LowerKeyTimestamp, timestamp)
	}
	createAPIKey(md)
	return metadata.NewOutgoingContext(ctx, md)
}

func GetErrorCode(err error) codes.Code {
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}

	var apiErr *xerror.APIError
	if errors.As(err, &apiErr) {
		return HTTPStatusToCode(apiErr.Code)
	}

	return codes.Unknown
}
