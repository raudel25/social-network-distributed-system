package socialnetwork

import (
	"path/filepath"
	"strings"

	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	follow_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_follow"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func userInFollowing(username string, followedUsername string) (bool, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")
	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	userFollowed, err := persistency.Load(node, path, &follow_pb.UserFollows{})
	if err != nil {
		return false, err
	}
	if userFollowed.FollowingUserIds == nil {
		return false, nil
	}
	for _, u := range userFollowed.FollowingUserIds {
		if u == followedUsername {
			return true, nil
		}
	}
	return false, nil
}

func follow(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")
	userFollows := &follow_pb.UserFollows{
		FollowingUserIds: make([]string, 0),
	}
	var err error
	exists, err := persistency.FileExists(node, path)
	if err != nil {
		return err
	}
	if exists {
		userFollows, err = persistency.Load(node, path, &follow_pb.UserFollows{})
		if err != nil {
			return err
		}
		if userFollows.FollowingUserIds == nil {
			userFollows.FollowingUserIds = make([]string, 0)
		}
	}
	userFollows.FollowingUserIds = append(userFollows.FollowingUserIds, otherUsername)
	persistency.Save(node, userFollows, path)
	return nil
}

func unfollow(username string, otherUsername string) error {
	path := filepath.Join("User", strings.ToLower(username), "Follow")

	exists, err := persistency.FileExists(node, path)

	if err != nil {
		return err
	}

	if !exists {
		return status.Errorf(404, "User %s following list not found", username)
	}

	userFollows := &follow_pb.UserFollows{
		FollowingUserIds: make([]string, 0),
	}

	if exists {
		userFollows, err = persistency.Load(node, path, &follow_pb.UserFollows{})
		if err != nil {
			return err
		}
	}
	for i, u := range userFollows.FollowingUserIds {
		if u == otherUsername {
			userFollows.FollowingUserIds = append(userFollows.FollowingUserIds[:i], userFollows.FollowingUserIds[i+1:]...)
			break
		}
	}
	persistency.Save(node, userFollows, path)
	return nil
}

func loadUserFollowing(username string) ([]*users_pb.User, error) {
	path := filepath.Join("User", strings.ToLower(username), "Follow")
	exists, err := persistency.FileExists(node, path)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, status.Errorf(codes.NotFound, "User %s following list not found", username)
	}

	userFollows, err := persistency.Load(node, path, &follow_pb.UserFollows{})

	if err != nil {
		return nil, err
	}

	users := make([]*users_pb.User, 0)

	for _, userId := range userFollows.FollowingUserIds {
		user, err := loadUser(userId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
