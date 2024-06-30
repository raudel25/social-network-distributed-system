package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
)

func existsUser(username string) bool {
	path := filepath.Join("User", strings.ToLower(username))
	return persistency.FileExists(node, path)
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
