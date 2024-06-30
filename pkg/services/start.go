package socialnetwork;

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
	node = chord.NewNode(&chord.Configuration{JoinAddress: "", SuccessorsSize: 5, HashSize: 4}, nil)
	node.Start("50050")
	go StartUserService(network, "0.0.0.0:50051")
	go StartAuthServer(network, "0.0.0.0:50052")
}
