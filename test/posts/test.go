package main

import (
	"time"

	"github.com/raudel25/social-network-distributed-system/pkg/logging"
	socialnetwork "github.com/raudel25/social-network-distributed-system/internal/services"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	tests "github.com/raudel25/social-network-distributed-system/test"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var token string

func main() {
	logging.SettingLogger(log.DebugLevel, ".")
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

	auth_client := socialnetwork_pb.NewAuthClient(auth_conn)
	post_cliet := socialnetwork_pb.NewPostServiceClient(posts_conn)

	password, _ := tests.HashPassword("hashedpassword")

	user := &socialnetwork_pb.User{
		Username:     "anabel",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}

	tests.TestSignUp(auth_client, user)
	token, _ = tests.TestLogin(auth_client, user.Username, "hashedpassword")
	tests.TestRepost(post_cliet, user.Username, "123", token)                          // not found
	tests.TestGetUserPosts(post_cliet, user.Username, token)                           // []
	postId := tests.TestCreatePost(post_cliet, "anabel", "This is a test post", token) // ok
	tests.TestGetPost(post_cliet, postId, token)                                       // post1
	tests.TestRepost(post_cliet, user.Username, postId, token)                         // ok
	tests.TestCreatePost(post_cliet, user.Username, "This is a test post 2", token)    // ok
	tests.TestGetUserPosts(post_cliet, user.Username, token)                           // [post1, repost1, post2]
}
