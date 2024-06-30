package socialnetwork

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// This interceptor is used for unary (single request-response) gRPC methods.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	_, err := ValidateRequest(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// This interceptor is used for streaming gRPC methods.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	_, err := ValidateRequest(ctx)
	if err != nil {
		return err
	}

	return handler(srv, ss)
}

// Is an interceptor for logging unary (single request-response) gRPC requests
func UnaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Get the peer address from the context
	p, _ := peer.FromContext(ctx)

	// Log the received request
	log.Printf("Request received - Method:%s From:%s", info.FullMethod, p.Addr.String())

	// Record the start time of the request processing
	start := time.Now()

	// Pass the request to the actual handler function
	h, err := handler(ctx, req)

	// Log the completed request
	log.Printf("Request completed - Method:%s\tDuration:%s\tError:%v",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}

// Is an interceptor for logging streaming gRPC requests
func StreamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Get the context from the server stream
	ctx := ss.Context()

	// Get the peer address from the context
	p, _ := peer.FromContext(ctx)

	// Log the received streaming request
	log.Printf("Streaming request received - Method:%s From:%s", info.FullMethod, p.Addr.String())

	// Record the start time of the request processing
	start := time.Now()

	// Pass the request to the actual handler function
	err := handler(srv, ss)

	// Log the completed streaming request
	log.Printf("Streaming Request completed - Method:%s\tDuration:%s\tError:%v",
		info.FullMethod,
		time.Since(start),
		err)

	return err
}
