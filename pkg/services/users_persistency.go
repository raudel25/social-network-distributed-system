package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func existsUser(username string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username))
	return persistency.FileExists(node, path)
}

func checkUsersExist(usernames ...string) error {
	for _, username := range usernames {
		exists, err := existsUser(username)
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to check user %s: %v", username, err)
		}
		if !exists {
			return status.Errorf(codes.NotFound, "User %s not found", username)
		}
	}
	return nil
}

func loadUser(username string) (*users_pb.User, error) {
	path := filepath.Join("User", strings.ToLower(username))
	user, err := persistency.Load(node, path, &users_pb.User{})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func saveUser(user *users_pb.User) error {
	path := filepath.Join("User", strings.ToLower(user.Username))
	return persistency.Save(node, user, path)
}
