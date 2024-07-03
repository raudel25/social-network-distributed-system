package main

import (
	"time"

	socialnetwork "github.com/raudel25/social-network-distributed-system/internal/services"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"github.com/raudel25/social-network-distributed-system/pkg/logging"
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

	follow_conn, err := grpc.NewClient("0.0.0.0:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	auth_client := socialnetwork_pb.NewAuthClient(auth_conn)
	follow_client := socialnetwork_pb.NewFollowServiceClient(follow_conn)

	password, _ := tests.HashPassword("hashedpassword")

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

	tests.TestSignUp(auth_client, user1)
	tests.TestSignUp(auth_client, user2)

	for {
		log.Println("========================================================= Start Test ==================================================================")

		token, _ = tests.TestLogin(auth_client, user1.Username, "hashedpassword")

		tests.TestGetFollowingUsers(follow_client, user1.Username, token) // []
		tests.TestGetFollowingUsers(follow_client, user2.Username, token) // []

		tests.TestFollow(follow_client, user1.Username, user2.Username, token) // ok

		time.Sleep(10 * time.Second) // tumbar el nodo de los datos

		tests.TestGetFollowingUsers(follow_client, user1.Username, token)        // [user2]
		tests.TestUnfollow(follow_client, user1.Username, user2.Username, token) // ok

		time.Sleep(10 * time.Second) // tumbar el nodo de los datos

		tests.TestGetFollowingUsers(follow_client, user1.Username, token)        // []
		tests.TestUnfollow(follow_client, user1.Username, user2.Username, token) // not following

		time.Sleep(15 * time.Second)
	}
}
