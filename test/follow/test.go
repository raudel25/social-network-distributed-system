package main

import (
	"context"
	"log"
	"time"

	socialnetwork "github.com/raudel25/social-network-distributed-system/pkg/services"
	auth_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_auth"
	follow_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_follow"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
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

	follow_conn, err := grpc.NewClient("0.0.0.0:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	auth_client := auth_pb.NewAuthClient(auth_conn)
	follow_client := follow_pb.NewFollowServiceClient(follow_conn)

	password, _ := hashPassword("hashedpassword")

	user1 := &users_pb.User{
		Username:     "anabel",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}

	user2 := &users_pb.User{
		Username:     "adriana",
		Name:         "jsdsjdb",
		PasswordHash: password,
		Email:        "tsddj@jdbj.com",
	}

	testSignUp(auth_client, user1)
	testLogin(auth_client, user1.Username, "hashedpassword")
	testSignUp(auth_client, user2)
	testLogin(auth_client, user2.Username, "hashedpassword")
	
	testGetFollowingUsers(follow_client, user1.Username)
	testFollow(follow_client, user1.Username, user2.Username)
	testFollow(follow_client, user1.Username, user2.Username)
	testGetFollowingUsers(follow_client, user1.Username)
	testFollow(follow_client, user2.Username, user1.Username)
	testGetFollowingUsers(follow_client, user2.Username)
}

func testSignUp(client auth_pb.AuthClient, user *users_pb.User) {
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

func testFollow(client follow_pb.FollowServiceClient, follower string, followed string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.FollowUser(ctx, &follow_pb.FollowUserRequest{
		UserId:       follower,
		TargetUserId: followed,
	})

	if err != nil {
		log.Printf("Error following user: %v", err)
	} else {
		log.Printf("User followed successfully")
	}
}

func testUnfollow(client follow_pb.FollowServiceClient, follower string, followed string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.UnfollowUser(ctx, &follow_pb.UnfollowUserRequest{
		UserId:       follower,
		TargetUserId: followed,
	})

	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
	} else {
		log.Printf("User unfollowed successfully")
	}
}

func testGetFollowingUsers(client follow_pb.FollowServiceClient, username string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetFollowing(ctx, &follow_pb.GetFollowingRequest{
		UserId: username,
	})

	if err != nil {
		log.Printf("Error getting following users: %v", err)
	} else {
		log.Printf("Following users: %v", response.GetFollowing())
	}
}
