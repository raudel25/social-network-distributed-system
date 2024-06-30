package main

import (
	"context"
	"log"
	"time"

	socialnetwork "github.com/raudel25/social-network-distributed-system/pkg/services"
	auth_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_auth"
	db_models_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_db"
	posts_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_posts"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var token string

func main() {
	rsaPrivateKeyPath := "pv.pem"
	rsaPublicteKeyPath := "pub.pem"
	network := "tcp"

	socialnetwork.Start(rsaPrivateKeyPath, rsaPublicteKeyPath, network)

	// Give the service a moment to start
	time.Sleep(time.Second)

	// Set up a connection to the server
	auth_conn, err := grpc.NewClient("0.0.0.0:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer auth_conn.Close()

	posts_conn, err := grpc.NewClient("0.0.0.0:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	auth_client := auth_pb.NewAuthClient(auth_conn)
	post_cliet := posts_pb.NewPostServiceClient(posts_conn)

	password, _ := hashPassword("hashedpassword")

	user := &db_models_pb.User{
		Username:     "anabel",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}

	testSignUp(auth_client, user)
	testLogin(auth_client, user.Username, "hashedpassword")
	testRepost(post_cliet, user.Username, "123")                          // not found
	testGetUserPosts(post_cliet, user.Username)                           // []
	postId := testCreatePost(post_cliet, "anabel", "This is a test post") // ok
	testGetPost(post_cliet, postId)                                       // post1
	testRepost(post_cliet, user.Username, postId)                         // ok
	testCreatePost(post_cliet, user.Username, "This is a test post 2")    // ok
	testGetUserPosts(post_cliet, user.Username)                           // [post1, repost1, post2]
}

func testSignUp(client auth_pb.AuthClient, user *db_models_pb.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.SignUp(ctx, &auth_pb.SignUpRequest{User: user})
	if err != nil {
		log.Printf("Error signing up: %v", err)
	} else {
		log.Println("Sign up successful")
	}
}

func testLogin(client auth_pb.AuthClient, username, password string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Login(ctx, &auth_pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Printf("Error logging in: %v", err)
	} else {
		token = response.GetToken()
		log.Printf("Login successful. Token: %s", token)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func testCreatePost(client posts_pb.PostServiceClient, username string, content string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.CreatePost(ctx, &posts_pb.CreatePostRequest{
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

func testGetPost(client posts_pb.PostServiceClient, postId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	post, err := client.GetPost(ctx, &posts_pb.GetPostRequest{
		PostId: postId,
	})

	if err != nil {
		log.Printf("Error getting post: %v", err)
	} else {
		log.Printf("Post retrieved successfully %v", post)
	}
}

func testRepost(client posts_pb.PostServiceClient, username string, postId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.Repost(ctx, &posts_pb.RepostRequest{
		UserId:         username,
		OriginalPostId: postId,
	})

	if err != nil {
		log.Printf("Error creating post: %v", err)
	} else {
		log.Printf("Post created successfully")
	}
}

func testGetUserPosts(client posts_pb.PostServiceClient, username string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetUserPosts(ctx, &posts_pb.GetUserPostsRequest{
		UserId: username,
	})

	if err != nil {
		log.Printf("Error getting user posts: %v", err)
	} else {
		log.Printf("Got user posts: %v", response.Posts)
	}
}
