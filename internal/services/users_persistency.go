package socialnetwork

import (
	"path/filepath"
	"strings"

	db_models "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func existsUser(username string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username))
	return persistency.FileExists(node, path)
}

func loadUser(username string) (*db_models.User, error) {
	path := filepath.Join("User", strings.ToLower(username))

	user, err := persistency.Load(node, path, &db_models.User{})

	if checkNotFound(err) {
		return nil, status.Errorf(codes.NotFound, "User %s not found", username)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load user %s: %v", username, err)
	}

	return user, nil
}

func saveUser(user *db_models.User) error {
	path := filepath.Join("User", strings.ToLower(user.Username))

	err := persistency.Save(node, user, path)

	if err != nil {
		return status.Errorf(codes.Internal, "Failed to save user: %v", err)
	}

	return nil
}
