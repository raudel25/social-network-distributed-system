package socialnetwork

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	posts_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_posts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServer struct {
	posts_pb.UnimplementedPostServiceServer
}

func (*PostServer) GetPost(_ context.Context, request *posts_pb.GetPostRequest) (*posts_pb.GetPostResponse, error) {
	postId := request.GetPostId()

	if !existsPost(postId) {
		return nil, status.Errorf(codes.NotFound, "Post %s not found", postId)
	}

	post, err := loadPost(postId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load post: %v", err)
	}

	return &posts_pb.GetPostResponse{Post: post}, nil
}

func (*PostServer) CreatePost(ctx context.Context, request *posts_pb.CreatePostRequest) (*posts_pb.CreatePostResponse, error) {
	if err := checkPermission(ctx, request.GetUserId()); err != nil {
		return nil, err
	}

	postID := fmt.Sprintf("%d", time.Now().UnixNano())

	post := &posts_pb.Post{
		PostId:    postID,
		UserId:    request.GetUserId(),
		Content:   request.GetContent(),
		Timestamp: time.Now().Unix(),
	}

	if err := savePost(post); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save post: %v", err)
	}

	if err := createUserPost(postID, request.GetUserId()); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save post to user: %v", err)
	}

	return &posts_pb.CreatePostResponse{Post: post}, nil
}

func (*PostServer) Repost(ctx context.Context, request *posts_pb.RepostRequest) (*posts_pb.RepostResponse, error) {
	if err := checkPermission(ctx, request.GetUserId()); err != nil {
		return nil, err
	}

	if !existsPost(request.OriginalPostId) {
		return nil, status.Errorf(codes.NotFound, "Post %s not found", request.OriginalPostId)
	}

	postID := fmt.Sprintf("%d", time.Now().UnixNano())

	post := &posts_pb.Post{
		PostId:         postID,
		UserId:         request.GetUserId(),
		Timestamp:      time.Now().Unix(),
		OriginalPostId: request.OriginalPostId,
	}

	if err := savePost(post); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save post: %v", err)
	}

	if err := createUserPost(postID, request.GetUserId()); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save post to user: %v", err)
	}

	return &posts_pb.RepostResponse{Post: post}, nil
}

func (*PostServer) GetUserPosts(_ context.Context, request *posts_pb.GetUserPostsRequest) (*posts_pb.GetUserPostsResponse, error) {
	userId := request.GetUserId()

	if !existsUser(userId) {
		return nil, status.Errorf(codes.NotFound, "User %s not found", userId)
	}

	posts, err := loadUserPosts(userId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load user posts: %v", err)
	}

	return &posts_pb.GetUserPostsResponse{Posts: posts}, nil
}

func StartPostsService(network string, address string) {
	log.Println("Post Service Started")

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

	posts_pb.RegisterPostServiceServer(s, &PostServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
