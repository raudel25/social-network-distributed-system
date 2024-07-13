package socialnetwork

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServer struct {
	socialnetwork_pb.UnimplementedPostServiceServer
}

func (*PostServer) GetPost(_ context.Context, request *socialnetwork_pb.GetPostRequest) (*socialnetwork_pb.GetPostResponse, error) {
	postId := request.GetPostId()

	post, err := loadPost(postId)

	if err != nil {
		return nil, err
	}

	return &socialnetwork_pb.GetPostResponse{Post: post}, nil
}

func (*PostServer) CreatePost(ctx context.Context, request *socialnetwork_pb.CreatePostRequest) (*socialnetwork_pb.CreatePostResponse, error) {
	if err := checkPermission(ctx, request.GetUserId()); err != nil {
		return nil, err
	}

	if len(request.GetContent()) > 140 {
		return nil, status.Errorf(codes.InvalidArgument, "Post content is too long")
	}

	if len(request.GetContent()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Posts must have content")
	}

	postID := fmt.Sprintf("%d", time.Now().UnixNano())

	post := &socialnetwork_pb.Post{
		PostId:    postID,
		UserId:    request.GetUserId(),
		Content:   request.GetContent(),
		Timestamp: time.Now().Unix(),
	}

	if err := savePost(post); err != nil {
		return nil, err
	}

	if err := addToPostsList(postID, request.GetUserId()); err != nil {
		return nil, err
	}

	return &socialnetwork_pb.CreatePostResponse{Post: post}, nil
}

func (*PostServer) Repost(ctx context.Context, request *socialnetwork_pb.RepostRequest) (*socialnetwork_pb.RepostResponse, error) {
	if err := checkPermission(ctx, request.GetUserId()); err != nil {
		return nil, err
	}

	originalPost, err := loadPost(request.OriginalPostId)

	if err != nil {
		return nil, err
	}

	postID := fmt.Sprintf("%d", time.Now().UnixNano())

	post := &socialnetwork_pb.Post{
		PostId:         postID,
		UserId:         request.GetUserId(),
		Content:        originalPost.Content,
		Timestamp:      time.Now().Unix(),
		OriginalPostId: request.GetOriginalPostId(),
	}

	if err := savePost(post); err != nil {
		return nil, err
	}

	if err := addToPostsList(postID, request.GetUserId()); err != nil {
		return nil, err
	}

	return &socialnetwork_pb.RepostResponse{Post: post}, nil
}

func (*PostServer) DeletePost(ctx context.Context, request *socialnetwork_pb.DeletePostRequest) (*socialnetwork_pb.DeletePostResponse, error) {
	postId := request.GetPostId()
	post, err := loadPost(postId)
	username := post.UserId

	if err != nil {
		return nil, err
	}

	if err := checkPermission(ctx, username); err != nil {
		return nil, err
	}

	if err := removePost(postId); err != nil {
		return nil, err
	}

	if err = removeFromPostsList(postId, username); err != nil {
		return nil, err
	}

	return &socialnetwork_pb.DeletePostResponse{}, nil
}

func (*PostServer) GetUserPosts(_ context.Context, request *socialnetwork_pb.GetUserPostsRequest) (*socialnetwork_pb.GetUserPostsResponse, error) {
	userId := request.GetUserId()

	if _, err := loadUser(userId); err != nil {
		return nil, err
	}

	posts, err := loadPostList(userId)

	if err != nil {
		return nil, err
	}

	return &socialnetwork_pb.GetUserPostsResponse{Posts: posts}, nil
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

	socialnetwork_pb.RegisterPostServiceServer(s, &PostServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ==============================================================================================================================================

func loadPostList(username string) ([]*socialnetwork_pb.Post, error) {
	posts := make([]*socialnetwork_pb.Post, 0)

	userPosts, err := loadPostsList(username)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load user posts: %v", err)
	}

	for _, postId := range userPosts.PostsIds {
		post, err := loadPost(postId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to load post %s in user posts: %v", postId, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
