package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	db_models_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ============================================== Post ==========================================================

func existsPost(postId string) (bool, error) {
	path := filepath.Join("Post", postId)
	return persistency.FileExists(node, path)
}

func checkPostsExist(postIds ...string) error {
	for _, postId := range postIds {
		exists, err := existsPost(postId)
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to check post %s: %v", postId, err)
		}
		if !exists {
			return status.Errorf(codes.NotFound, "Post %s not found", postId)
		}
	}
	return nil
}

func loadPost(postId string) (*db_models_pb.Post, error) {
	path := filepath.Join("Post", postId)
	post, err := persistency.Load(node, path, &db_models_pb.Post{})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func savePost(post *db_models_pb.Post) error {
	path := filepath.Join("Post", post.PostId)
	return persistency.Save(node, post, path)
}

// ========================================= User-Post relationship =====================================================

func addToPostsList(postId string, username string) error {
	path := filepath.Join("User", strings.ToLower(username), "Posts")
	posts := &db_models_pb.UserPosts{
		PostsIds: make([]string, 0),
	}
	var err error

	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return err
	}

	if exists {
		posts, err = persistency.Load(node, path, &db_models_pb.UserPosts{})
		if err != nil {
			return err
		}
	}

	posts.PostsIds = append(posts.PostsIds, postId)
	return persistency.Save(node, posts, path)
}

func loadPostsList(username string) ([]*db_models_pb.Post, error) {
	path := filepath.Join("User", strings.ToLower(username), "Posts")

	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return nil, err
	}

	posts := make([]*db_models_pb.Post, 0)

	if !exists {
		return posts, nil
	}

	userPosts, err := persistency.Load(node, path, &db_models_pb.UserPosts{})
	if err != nil {
		return nil, err
	}
	for _, postId := range userPosts.PostsIds {
		post, err := loadPost(postId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
