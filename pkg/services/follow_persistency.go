package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	db_models "github.com/raudel25/social-network-distributed-system/pkg/services/grpc"
	"google.golang.org/grpc/status"
)

// ============================================== User-User follow relationship ==========================================================

func existsInFollowingList(username string, otherUsername string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")
	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	userFollowed, err := persistency.Load(node, path, &db_models.UserFollows{})
	if err != nil {
		return false, err
	}
	if userFollowed.FollowingUserIds == nil {
		return false, nil
	}
	for _, u := range userFollowed.FollowingUserIds {
		if u == otherUsername {
			return true, nil
		}
	}
	return false, nil
}

func loadFollowingList(username string) ([]*db_models.User, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")
	exists, err := persistency.FileExists(node, path)

	if err != nil {
		return nil, err
	}

	users := make([]*db_models.User, 0)

	if !exists {
		return users, nil
	}

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})

	if err != nil {
		return nil, err
	}

	for _, userId := range userFollows.FollowingUserIds {
		user, err := loadUser(userId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func addToFollowingList(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows := &db_models.UserFollows{
		FollowingUserIds: make([]string, 0),
	}

	var err error

	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return err
	}

	if exists {
		userFollows, err = persistency.Load(node, path, &db_models.UserFollows{})
		if err != nil {
			return err
		}
	}

	userFollows.FollowingUserIds = append(userFollows.FollowingUserIds, otherUsername)
	return persistency.Save(node, userFollows, path)
}

func removeFromFollowingList(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	exists, err := persistency.FileExists(node, path)

	if err != nil {
		return err
	}

	if !exists {
		return status.Errorf(404, "User %s following list not found", username)
	}

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})
	if err != nil {
		return err
	}

	for i, u := range userFollows.FollowingUserIds {
		if u == otherUsername {
			userFollows.FollowingUserIds = append(userFollows.FollowingUserIds[:i], userFollows.FollowingUserIds[i+1:]...)
			break
		}
	}

	return persistency.Save(node, userFollows, path)
}
