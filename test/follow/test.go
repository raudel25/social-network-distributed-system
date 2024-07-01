package main

import (
	"context"
	"time"

	"github.com/raudel25/social-network-distributed-system/pkg/logging"
	socialnetwork "github.com/raudel25/social-network-distributed-system/pkg/services"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	follow_conn, err := grpc.NewClient("0.0.0.0:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	auth_client := socialnetwork_pb.NewAuthClient(auth_conn)
	follow_client := socialnetwork_pb.NewFollowServiceClient(follow_conn)

	password, _ := hashPassword("hashedpassword")

	user1 := &socialnetwork_pb.User{
		Username:     "anabel",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}

	user2 := &socialnetwork_pb.User{
		Username:     "adriana",
		Name:         "jsdsjdb",
		PasswordHash: password,
		Email:        "tsddj@jdbj.com",
	}

	testSignUp(auth_client, user1)
	testSignUp(auth_client, user2)
	testLogin(auth_client, user2.Username, "hashedpassword")

	testGetFollowingUsers(follow_client, user1.Username)      // []
	testFollow(follow_client, user1.Username, user2.Username) // not authorized

	testFollow(follow_client, user2.Username, user1.Username) // ok
	testFollow(follow_client, user2.Username, user1.Username) // already following
	testGetFollowingUsers(follow_client, user2.Username)      // [user1]

	testLogin(auth_client, user1.Username, "hashedpassword")
	testFollow(follow_client, user1.Username, user2.Username)   // ok
	testGetFollowingUsers(follow_client, user1.Username)        // [user2]
	testUnfollow(follow_client, user1.Username, user2.Username) // ok
	testGetFollowingUsers(follow_client, user1.Username)        // []
	testUnfollow(follow_client, user1.Username, user2.Username) // not following
}

func testSignUp(client socialnetwork_pb.AuthClient, user *socialnetwork_pb.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.SignUp(ctx, &socialnetwork_pb.SignUpRequest{User: user})
	if err != nil {
		log.Printf("Error signing up: %v", err)
	} else {
		log.Println("Sign up successful")
	}
}

func testLogin(client socialnetwork_pb.AuthClient, username, password string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Login(ctx, &socialnetwork_pb.LoginRequest{
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

func testFollow(client socialnetwork_pb.FollowServiceClient, follower string, followed string) {
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

func testUnfollow(client socialnetwork_pb.FollowServiceClient, follower string, followed string) {
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

func testGetFollowingUsers(client socialnetwork_pb.FollowServiceClient, username string) {
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
