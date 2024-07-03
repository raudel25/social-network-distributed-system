package tests

import (
	"context"
	"time"

	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

//==========================================================USERS====================================================================

func TestGetUser(client socialnetwork_pb.UserServiceClient, username string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetUser(ctx, &socialnetwork_pb.GetUserRequest{Username: username})
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		log.Printf("Got user: %v", response.GetUser())
	}
}

func TestEditUser(client socialnetwork_pb.UserServiceClient, user *socialnetwork_pb.User, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.EditUser(ctx, &socialnetwork_pb.EditUserRequest{User: user})
	if err != nil {
		log.Printf("Error editing user: %v", err)
	} else {
		log.Printf("User edited successfully")
	}

	TestGetUser(client, user.Username, token)
}

//==========================================================AUTH====================================================================

func TestSignUp(client socialnetwork_pb.AuthClient, user *socialnetwork_pb.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.SignUp(ctx, &socialnetwork_pb.SignUpRequest{User: user})
	if err != nil {
		log.Printf("Error signing up: %v", err)
	} else {
		log.Println("Sign up successful")
	}
}

func TestLogin(client socialnetwork_pb.AuthClient, username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Login(ctx, &socialnetwork_pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Printf("Error logging in: %v", err)
		return "", err
	} else {
		token := response.GetToken()
		log.Printf("Login successful. Token: %s", token)
		return token, nil
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

//==========================================================POSTS====================================================================

func TestCreatePost(client socialnetwork_pb.PostServiceClient, username string, content string, token string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.CreatePost(ctx, &socialnetwork_pb.CreatePostRequest{
		UserId:  username,
		Content: content,
	})

	if err != nil {
		log.Printf("Error creating post: %v", err)
	} else {
		log.Printf("Post created successfully %v", response.GetPost())
	}

	return response.GetPost().PostId
}

func TestGetPost(client socialnetwork_pb.PostServiceClient, postId string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	post, err := client.GetPost(ctx, &socialnetwork_pb.GetPostRequest{
		PostId: postId,
	})

	if err != nil {
		log.Printf("Error getting post: %v", err)
	} else {
		log.Printf("Post retrieved successfully %v", post)
	}
}

func TestRepost(client socialnetwork_pb.PostServiceClient, username string, postId string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.Repost(ctx, &socialnetwork_pb.RepostRequest{
		UserId:         username,
		OriginalPostId: postId,
	})

	if err != nil {
		log.Printf("Error creating post: %v", err)
	} else {
		log.Printf("Post created successfully")
	}
}

func TestGetUserPosts(client socialnetwork_pb.PostServiceClient, username string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetUserPosts(ctx, &socialnetwork_pb.GetUserPostsRequest{
		UserId: username,
	})

	if err != nil {
		log.Printf("Error getting user posts: %v", err)
	} else {
		log.Printf("Got user posts: %v", response.Posts)
	}
}

//==========================================================FOLLOW====================================================================

func TestFollow(client socialnetwork_pb.FollowServiceClient, follower string, followed string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.FollowUser(ctx, &socialnetwork_pb.FollowUserRequest{
		UserId:       follower,
		TargetUserId: followed,
	})

	if err != nil {
		log.Printf("Error following user: %v", err)
	} else {
		log.Printf("User followed successfully")
	}
}

func TestUnfollow(client socialnetwork_pb.FollowServiceClient, follower string, followed string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.UnfollowUser(ctx, &socialnetwork_pb.UnfollowUserRequest{
		UserId:       follower,
		TargetUserId: followed,
	})

	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
	} else {
		log.Printf("User unfollowed successfully")
	}
}

func TestGetFollowingUsers(client socialnetwork_pb.FollowServiceClient, username string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetFollowing(ctx, &socialnetwork_pb.GetFollowingRequest{
		UserId: username,
	})

	if err != nil {
		log.Printf("Error getting following users: %v", err)
	} else {
		log.Printf("Following users: %v", response.GetFollowing())
	}
}
