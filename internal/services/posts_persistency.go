package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	db_models "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ============================================== Post ==========================================================

func loadPost(postId string) (*db_models.Post, error) {
	path := filepath.Join("Post", postId)

	post, err := persistency.Load(node, path, &db_models.Post{})

	if checkNotFound(err) {
		return nil, status.Errorf(codes.NotFound, "Post %s not found", postId)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load post: %v", err)
	}

	return post, nil
}

func savePost(post *db_models.Post) error {
	path := filepath.Join("Post", post.PostId)

	err := persistency.Save(node, post, path)

	if err != nil {
		return status.Errorf(codes.Internal, "Failed to save post: %v", err)
	}

	return nil
}

// ========================================= User-Post relationship =====================================================

func addToPostsList(postId string, username string) error {
	path := filepath.Join("User", strings.ToLower(username), "Posts")

	posts, err := persistency.Load(node, path, &db_models.UserPosts{})

	if checkNotFound(err) {
		posts = &db_models.UserPosts{
			PostsIds: make([]string, 0),
		}
	} else if err != nil {
		return status.Errorf(codes.Internal, "Failed to load user posts: %v", err)
	}

	posts.PostsIds = append(posts.PostsIds, postId)

	err = persistency.Save(node, posts, path)

	if err != nil {
		return status.Errorf(codes.Internal, "Failed to save post to user: %v", err)
	}

	return nil
}

func loadPostsList(username string) ([]*db_models.Post, error) {
	path := filepath.Join("User", strings.ToLower(username), "Posts")

	posts := make([]*db_models.Post, 0)

	userPosts, err := persistency.Load(node, path, &db_models.UserPosts{})

	if checkNotFound(err) {
		return posts, nil
	}

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
