package main

import (
	"context"
	"log"
	"time"

	socialnetwork "github.com/raudel25/social-network-distributed-system/pkg/services"
	auth_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_auth"
	db_models_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_db"
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

	// Create a client
	auth_client := auth_pb.NewAuthClient(auth_conn)

	users_conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer auth_conn.Close()
	users_client := users_pb.NewUserServiceClient(users_conn)

	password, _ := hashPassword("hashedpassword")

	user1 := &db_models_pb.User{
		Username:     "hola",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}

	user2 := &db_models_pb.User{
		Username:     "anabel",
		Name:         "Test User",
		PasswordHash: password,
		Email:        "testuser@example.com",
	}
	editedUser1 := &db_models_pb.User{
		Username: "hola",
		Name:     "Anabel Benítez",
		Email:    "testuser@example.com",
	}

	testSignUp(auth_client, user1)                           // ok
	testSignUp(auth_client, user2)                           // ok
	testLogin(auth_client, user2.Username, "hashedpassword") // valid token
	testGetUser(users_client, user1.Username)                // user1
	testEditUser(users_client, editedUser1)                  // not authorized
}

func testGetUser(client users_pb.UserServiceClient, username string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := client.GetUser(ctx, &users_pb.GetUserRequest{Username: username})
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		log.Printf("Got user: %v", response.GetUser())
	}
}

func testEditUser(client users_pb.UserServiceClient, user *db_models_pb.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.EditUser(ctx, &users_pb.EditUserRequest{User: user})
	if err != nil {
		log.Printf("Error editing user: %v", err)
	} else {
		log.Printf("User edited successfully")
	}

	testGetUser(client, user.Username)
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
