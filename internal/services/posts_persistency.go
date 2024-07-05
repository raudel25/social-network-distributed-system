package socialnetwork

import (
	"path/filepath"
	"strings"

	db_models "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
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

func loadPostDto(postId string) (*db_models.PostDto, error) {
	post, err := loadPost(postId)

	if err != nil {
		return nil, err
	}

	postResponse := postToPostDto(post, db_models.PostType_ORIGINAL)

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

func loadOriginalPost(postId string, times int) (*db_models.PostDto, db_models.PostType, error) {
	originalPost, err := loadPost(postId)

	postType := db_models.PostType_REPOST

	if times > 1 {
		postType = db_models.PostType_REPOST_OF_REPOST
	}

	if err != nil {
		return nil, postType, err
	}

	if originalPost.OriginalPostId == "" {
		responsePost := postToPostDto(originalPost, db_models.PostType_ORIGINAL)
		return responsePost, postType, nil
	}

	return loadOriginalPost(originalPost.OriginalPostId, times+1)
}

func postToPostDto(post *db_models.Post, postType db_models.PostType) *db_models.PostDto {
	return &db_models.PostDto{
		PostId:    post.PostId,
		UserId:    post.UserId,
		Content:   post.Content,
		Timestamp: post.Timestamp,
		PostType:  postType,
	}
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

func loadPostsList(username string) ([]*db_models.PostDto, error) {
	path := filepath.Join("User", strings.ToLower(username), "Posts")

	posts := make([]*db_models.PostDto, 0)

	userPosts, err := persistency.Load(node, path, &db_models.UserPosts{})

	if checkNotFound(err) {
		return posts, nil
	}

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
