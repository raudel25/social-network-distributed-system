package socialnetwork

import (
	"path/filepath"
	"strings"

	db_models "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ============================================== User-User follow relationship ==========================================================

func existsInFollowingList(username string, otherUsername string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollowed, err := persistency.Load(node, path, &db_models.UserFollows{})

	if checkNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, status.Errorf(codes.Internal, "Failed to load following %v", err)
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

	users := make([]*db_models.User, 0)

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})

	if checkNotFound(err) {
		return users, nil
	}

	if err != nil {
		return nil, err
	}

	for _, userId := range userFollows.FollowingUserIds {
		user, err := loadUser(userId)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = ""
		users = append(users, user)
	}

	return users, nil
}

func addToFollowingList(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})

	if checkNotFound(err) {
		userFollows = &db_models.UserFollows{
			FollowingUserIds: make([]string, 0),
		}
	} else if err != nil {
		return err
	}

	userFollows.FollowingUserIds = append(userFollows.FollowingUserIds, otherUsername)
	return persistency.Save(node, userFollows, path)
}

func removeFromFollowingList(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	userFollows, err := persistency.Load(node, path, &db_models.UserFollows{})

	if checkNotFound(err) {
		return status.Errorf(codes.NotFound, "User %s following list not found", username)
	}

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
