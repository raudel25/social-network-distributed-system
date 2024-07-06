package socialnetwork

import (
	log "github.com/sirupsen/logrus"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
)

var (
	node       *chord.Node
	rsaPrivate string
	rsaPublic  string
)

func Start(rsaPrivateKeyPath string, rsaPublicteKeyPath string, network string) {
	var err error

	rsaPrivate = rsaPrivateKeyPath
	rsaPublic = rsaPublicteKeyPath

	node, err = chord.DefaultNode()

	if err != nil {
		log.Fatalf("Can't start chord node")
	}

	port := "50050"
	broadListen := "6000"
	broadRequest := "7000"
	node.Start(port, broadListen, broadRequest)

	go StartUserService(network, "0.0.0.0:50051")
	go StartAuthServer(network, "0.0.0.0:50052")
	go StartPostsService(network, "0.0.0.0:50053")
	go StartFollowService(network, "0.0.0.0:50054")
}
