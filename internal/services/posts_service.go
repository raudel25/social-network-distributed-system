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

	post, err := loadPostDto(postId)

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

	return &socialnetwork_pb.CreatePostResponse{Post: postToPostDto(post, 0)}, nil
}

func (*PostServer) Repost(ctx context.Context, request *socialnetwork_pb.RepostRequest) (*socialnetwork_pb.RepostResponse, error) {
	if err := checkPermission(ctx, request.GetUserId()); err != nil {
		return nil, err
	}

	_, err := loadPost(request.OriginalPostId)

	if err != nil {
		return nil, err
	}

	postID := fmt.Sprintf("%d", time.Now().UnixNano())

	post := &socialnetwork_pb.Post{
		PostId:         postID,
		UserId:         request.GetUserId(),
		Content:        request.GetContent(),
		Timestamp:      time.Now().Unix(),
		OriginalPostId: request.GetOriginalPostId(),
	}

	if err := savePost(post); err != nil {
		return nil, err
	}

	if err := addToPostsList(postID, request.GetUserId()); err != nil {
		return nil, err
	}

	postResponse, err := loadPostDto(postID)

	if err != nil {
		return nil, err
	}

	return &socialnetwork_pb.RepostResponse{Post: postResponse}, nil
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

	posts, err := loadPostDtosList(userId)

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

func loadPostDtosList(username string) ([]*socialnetwork_pb.PostDto, error) {
	posts := make([]*socialnetwork_pb.PostDto, 0)

	userPosts, err := loadPostsList(username)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load user posts: %v", err)
	}

	for _, postId := range userPosts.PostsIds {
		post, err := loadPostDto(postId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to load post %s in user posts: %v", postId, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func loadPostDto(postId string) (*socialnetwork_pb.PostDto, error) {
	post, err := loadPost(postId)

	if err != nil {
		return nil, err
	}

	postResponse := postToPostDto(post, socialnetwork_pb.PostType_ORIGINAL)

	if post.OriginalPostId != "" {
		postResponseOriginal, postType, err := loadOriginalPost(postId, 0)

		// asumimos que un not found es que borraron algun post en el camino
		if err != nil && !checkNotFound(err) {
			return nil, err
		}

		postResponse.PostType = postType
		postResponse.OriginalPost = postResponseOriginal
	}

	return postResponse, nil
}

func loadOriginalPost(postId string, times int) (*socialnetwork_pb.PostDto, socialnetwork_pb.PostType, error) {
	originalPost, err := loadPost(postId)

	postType := socialnetwork_pb.PostType_REPOST

	if times > 1 {
		postType = socialnetwork_pb.PostType_REPOST_OF_REPOST
	}

	if err != nil {
		return nil, postType, err
	}

	if originalPost.OriginalPostId == "" {
		responsePost := postToPostDto(originalPost, socialnetwork_pb.PostType_ORIGINAL)
		return responsePost, postType, nil
	}

	return loadOriginalPost(originalPost.OriginalPostId, times+1)
}

func postToPostDto(post *socialnetwork_pb.Post, postType socialnetwork_pb.PostType) *socialnetwork_pb.PostDto {
	return &socialnetwork_pb.PostDto{
		PostId:    post.PostId,
		UserId:    post.UserId,
		Content:   post.Content,
		Timestamp: post.Timestamp,
		PostType:  postType,
	}
}
