package socialnetwork

import (
	"path/filepath"
	"strings"

	db_models "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
)

// ============================================== User-User follow relationship ==========================================================

func loadFollowingList(username string) (*db_models.UserFollows, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})

	if checkNotFound(err) {
		return &db_models.UserFollows{FollowingUserIds: make([]string, 0)}, nil
	}

	if err != nil {
		return nil, err
	}

	return userFollows, nil
}

func addToFollowingList(username string, otherUsername string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows, err := loadFollowingList(username)

	if err != nil {
		return false, err
	}

	ok := true
	for _, u := range userFollows.FollowingUserIds {
		if u == otherUsername {
			ok = false
		}
	}

	if ok {
		userFollows.FollowingUserIds = append(userFollows.FollowingUserIds, otherUsername)
	}

	return ok, persistency.Save(node, userFollows, path)
}

func removeFromFollowingList(username string, otherUsername string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows, err := loadFollowingList(username)

	if err != nil {
		return false, err
	}

	ok := false

	for i, u := range userFollows.FollowingUserIds {
		if u == otherUsername {
			userFollows.FollowingUserIds = append(userFollows.FollowingUserIds[:i], userFollows.FollowingUserIds[i+1:]...)
			ok = true
		}
	}

	return ok, persistency.Save(node, userFollows, path)
}
