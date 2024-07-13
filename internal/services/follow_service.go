package socialnetwork

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FollowServer struct {
	*socialnetwork_pb.UnimplementedFollowServiceServer
}

func (*FollowServer) FollowUser(ctx context.Context, request *socialnetwork_pb.FollowUserRequest) (*socialnetwork_pb.FollowUserResponse, error) {
	username := request.GetUserId()
	targetUsername := request.GetTargetUserId()

	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if username == targetUsername {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot follow yourself")
	}

	if _, err := loadUser(username); err != nil {
		return nil, err
	}

	if _, err := loadUser(targetUsername); err != nil {
		return nil, err
	}

	ok, err := addToFollowingList(username, targetUsername)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to follow user %v", err)
	}

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Already following user %s", targetUsername)
	}

	return &socialnetwork_pb.FollowUserResponse{}, nil
}

func (*FollowServer) UnfollowUser(ctx context.Context, request *socialnetwork_pb.UnfollowUserRequest) (*socialnetwork_pb.UnfollowUserResponse, error) {
	username := request.GetUserId()
	targetUsername := request.GetTargetUserId()

	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if username == targetUsername {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot unfollow yourself")
	}

	if _, err := loadUser(username); err != nil {
		return nil, err
	}

	if _, err := loadUser(targetUsername); err != nil {
		return nil, err
	}

	ok, err := removeFromFollowingList(username, targetUsername)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unfollow user %v", err)
	}

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Not following user %s", targetUsername)
	}

	return &socialnetwork_pb.UnfollowUserResponse{}, nil
}

func (*FollowServer) GetFollowing(ctx context.Context, request *socialnetwork_pb.GetFollowingRequest) (*socialnetwork_pb.GetFollowingResponse, error) {
	username := request.GetUserId()

	if _, err := loadUser(username); err != nil {
		return nil, err
	}

	following, err := loadFollowingUsenamesList(username)

	if err != nil {
		return nil, err
	}

	return &socialnetwork_pb.GetFollowingResponse{FollowingUsernames: following}, nil
}

func StartFollowService(network, address string) {
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

	socialnetwork_pb.RegisterFollowServiceServer(s, &FollowServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadFollowingUsenamesList(username string) ([]string, error) {
	users := make([]string, 0)

	userFollows, err := loadFollowingList(username)

	if err != nil {
		return nil, err
	}

	users = append(users, userFollows.FollowingUserIds...)

	return users, nil
}
