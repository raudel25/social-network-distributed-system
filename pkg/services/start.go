package socialnetwork

import (
	"github.com/raudel25/social-network-distributed-system/pkg/chord"
)

var (
	node       *chord.Node
	rsaPrivate string
	rsaPublic  string
)

func Start(rsaPrivateKeyPath string, rsaPublicteKeyPath string, network string) {
	rsaPrivate = rsaPrivateKeyPath
	rsaPublic = rsaPublicteKeyPath
	node = chord.NewNode(&chord.Configuration{SuccessorsSize: 5, HashSize: 4}, nil)
	node.Start("5000", "8000")
	go StartUserService(network, "0.0.0.0:50051")
	go StartAuthServer(network, "0.0.0.0:50052")
	go StartPostsService(network, "0.0.0.0:50053")
	go StartFollowService(network, "0.0.0.0:50054")
}
