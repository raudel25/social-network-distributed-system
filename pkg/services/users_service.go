package socialnetwork

import (
	"context"
	"log"
	"net"
	"path/filepath"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	users_pb.UnimplementedUserServiceServer
}

func (*UserServer) GetUser(_ context.Context, request *users_pb.GetUserRequest) (*users_pb.GetUserResponse, error) {
	user, err := loadUser(request.GetUsername())

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

	if err := saveUser(request.GetUser()); err != nil {
		return nil, err
	}

	return &users_pb.EditUserResponse{}, nil
}

func loadUser(username string) (*users_pb.User, error) {
	user := &users_pb.User{}
	path := filepath.Join("User", strings.ToLower(username))
	return persistency.Load(node, path, user)
}

func saveUser(user *users_pb.User) error {
	user.Username = strings.ToLower(user.Username)
	path := filepath.Join("User", user.Username)
	if persistency.FileExists(node, path) {
		return status.Error(codes.AlreadyExists, "Username is taken")
	}
	return persistency.Save(node, user, path)
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
