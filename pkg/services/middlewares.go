package socialnetwork

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// This interceptor is used for unary (single request-response) gRPC methods.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	_, err := validateRequest(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// This interceptor is used for streaming gRPC methods.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	_, err := validateRequest(ctx)
	if err != nil {
		return err
	}

	return handler(srv, ss)
}

// This interceptor is used for logging unary (single request-response) gRPC requests
func UnaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, _ := peer.FromContext(ctx)

	log.Infof("Request received - Method:%s From:%s", info.FullMethod, p.Addr.String())

	start := time.Now()

	h, err := handler(ctx, req)

	log.Infof("Request completed - Method:%s\tDuration:%s\tError:%v",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}

// This interceptor is used for logging streaming gRPC requests
func StreamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	p, _ := peer.FromContext(ctx)

	log.Infof("Streaming request received - Method:%s From:%s", info.FullMethod, p.Addr.String())

	start := time.Now()

	err := handler(srv, ss)

	log.Infof("Streaming Request completed - Method:%s\tDuration:%s\tError:%v",
		info.FullMethod,
		time.Since(start),
		err)

	return err
}
