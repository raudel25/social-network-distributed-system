package socialnetwork;

import (
	"context"
	"log"
	"net"
	"path/filepath"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/users"
	"google.golang.org/grpc"
)

type UserServer struct {
	users_pb.UnimplementedUserServiceServer
}

func (*UserServer) GetUser(_ context.Context, request *users_pb.GetUserRequest) (*users_pb.GetUserResponse, error) {
	username := request.GetUsername()

	user := &users_pb.User{}
	user, err := persistency.Load(node, filepath.Join("User", username), user)

	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &users_pb.GetUserResponse{User: user}, nil
}

func (s *UserServer) EditUser(ctx context.Context, request *users_pb.EditUserRequest) (*users_pb.EditUserResponse, error) {

	if err := checkPermission(ctx, request.GetUser().Username); err != nil {
		return nil, err
	}

	err := persistency.Save(node, request.GetUser(), filepath.Join("User", request.GetUser().Username))
	if err != nil {
		return nil, err
	}

	return &users_pb.EditUserResponse{}, nil
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

	users_pb.RegisterUserServiceServer(s, &UserServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
