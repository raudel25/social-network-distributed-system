package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	posts_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_posts"
)

func existsPost(postId string) bool {
	path := filepath.Join("Post", postId)
	return persistency.FileExists(node, path)
}

func loadPost(postId string) (*posts_pb.Post, error) {
	path := filepath.Join("Post", postId)
	post, err := persistency.Load(node, path, &posts_pb.Post{})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func savePost(post *posts_pb.Post) error {
	path := filepath.Join("Post", post.PostId)
	return persistency.Save(node, post, path)
}

func createUserPost(postId string, username string) error {
	path := filepath.Join("User", strings.ToLower(username), "Posts")
	posts := &posts_pb.UserPosts{
		PostsIds: make([]string, 0),
	}
	var err error

	if persistency.FileExists(node, path) {
		posts, err = persistency.Load(node, path, &posts_pb.UserPosts{})
		if err != nil {
			return err
		}
		if posts.PostsIds == nil {
			posts.PostsIds = make([]string, 0)
		}
	}

	posts.PostsIds = append(posts.PostsIds, postId)
	return persistency.Save(node, posts, path)
}

func loadUserPosts(username string) ([]*posts_pb.Post, error) {
	path := filepath.Join("User", strings.ToLower(username), "Posts")
	userPosts, err := persistency.Load(node, path, &posts_pb.UserPosts{})
	if err != nil {
		return nil, err
	}
	posts := make([]*posts_pb.Post, 0)
	for _, postId := range userPosts.PostsIds {
		post, err := loadPost(postId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
