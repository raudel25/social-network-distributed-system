package socialnetwork

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"google.golang.org/grpc"
)

type UserServer struct {
	socialnetwork_pb.UnimplementedUserServiceServer
}

func (*UserServer) GetUser(_ context.Context, request *socialnetwork_pb.GetUserRequest) (*socialnetwork_pb.GetUserResponse, error) {
	username := request.GetUsername()

	user, err := loadUser(username)

	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &socialnetwork_pb.GetUserResponse{User: user}, nil
}

func (s *UserServer) EditUser(ctx context.Context, request *socialnetwork_pb.EditUserRequest) (*socialnetwork_pb.EditUserResponse, error) {
	username := request.GetUser().Username
	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if _, err := loadUser(username); err != nil {
		return nil, err
	}

	if err := saveUser(request.GetUser()); err != nil {
		return nil, err
	}

	return &socialnetwork_pb.EditUserResponse{}, nil
}

func StartUserService(network string, address string) {
	log.Println("User Service Started")

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				UnaryLoggingInterceptor,
				UnaryServerInterceptor,
			),
		), grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				StreamLoggingInterceptor,
				StreamServerInterceptor,
			),
		),
	)

	socialnetwork_pb.RegisterUserServiceServer(s, &UserServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
