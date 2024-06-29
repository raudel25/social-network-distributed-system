package services

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// This interceptor is used for unary (single request-response) gRPC methods.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Validate the request
	token, err := ValidateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Extract metadata from the incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("Error extracting metadata from context")
		return nil, status.Error(codes.Internal, "")
	}

	// Append the username from the JWT token claims to the metadata
	md.Append("username", token.Claims.(jwt.MapClaims)["sub"].(string))

	// Create a new incoming context with the updated metadata
	newCtx := metadata.NewIncomingContext(ctx, md)

	// Pass the request to the actual handler function with the new context
	return handler(newCtx, req)
}

// This interceptor is used for streaming gRPC methods.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Get the context from the server stream
	ctx := ss.Context()

	// Validate the request
	_, err := ValidateRequest(ctx)
	if err != nil {
		return err
	}

	// Pass the request to the actual handler function
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
