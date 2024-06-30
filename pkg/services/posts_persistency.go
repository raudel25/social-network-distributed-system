package socialnetwork

import (
	"path/filepath"

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

func createUserPost(postId string, userId string) error {
	path := filepath.Join("User", userId, "Posts")
	posts := &posts_pb.UserPosts{
		PostsIds: make([]string, 0),
	}
	var err error

	if persistency.FileExists(node, path) {
		posts, err = persistency.Load(node, path, &posts_pb.UserPosts{})
		if err != nil {
			return err
		}
	}

	posts.PostsIds = append(posts.PostsIds, postId)
	return persistency.Save(node, posts, path)
}

func loadUserPosts(userId string) ([]*posts_pb.Post, error) {
	path := filepath.Join("User", userId, "Posts")
	userPosts, err := persistency.Load(node, path, &posts_pb.UserPosts{})
	if err != nil {
		return nil, err
	}
	posts := make([]*posts_pb.Post, 0)
	for i := range userPosts.PostsIds {
		post, err := loadPost(userPosts.PostsIds[i])
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
