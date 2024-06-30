package socialnetwork

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	follow_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_follow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FollowServer struct {
	*follow_pb.UnimplementedFollowServiceServer
}

func (*FollowServer) FollowUser(ctx context.Context, request *follow_pb.FollowUserRequest) (*follow_pb.FollowUserResponse, error) {
	username := request.GetUserId()
	targetUsername := request.GetTargetUserId()

	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if username == targetUsername {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot follow yourself")
	}

	if err := checkUsersExist(username, targetUsername); err != nil {
		return nil, err
	}

	following, err := existsInFollowingList(username, targetUsername)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load following %v", err)
	}
	if following {
		return nil, status.Errorf(codes.AlreadyExists, "Already following user")
	}

	if err := addToUserFollowingList(username, targetUsername); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to follow user %v", err)
	}

	return &follow_pb.FollowUserResponse{}, nil
}

func (*FollowServer) UnfollowUser(ctx context.Context, request *follow_pb.UnfollowUserRequest) (*follow_pb.UnfollowUserResponse, error) {
	username := request.GetUserId()
	targetUsername := request.GetTargetUserId()

	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if username == targetUsername {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot unfollow yourself")
	}

	if err := checkUsersExist(username, targetUsername); err != nil {
		return nil, err
	}

	following, err := existsInFollowingList(username, targetUsername)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load following %v", err)
	}
	if !following {
		return nil, status.Errorf(codes.NotFound, "Not following user")
	}

	if err := removeFromUserFollowingList(username, targetUsername); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unfollow user %v", err)
	}

	return &follow_pb.UnfollowUserResponse{}, nil
}

func (*FollowServer) GetFollowing(ctx context.Context, request *follow_pb.GetFollowingRequest) (*follow_pb.GetFollowingResponse, error) {
	username := request.GetUserId()

	if err := checkUsersExist(username); err != nil {
		return nil, err
	}

	following, err := loadFollowing(username)

	if err != nil {
		return nil, err
	}

	return &follow_pb.GetFollowingResponse{Following: following}, nil
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

	follow_pb.RegisterFollowServiceServer(s, &FollowServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}