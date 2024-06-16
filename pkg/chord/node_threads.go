package chord

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func (n *Node) threadStabilize() {
	log.Println("Stabilize thread started")

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.stabilize()
		}
	}
}

func (n *Node) threadListen(s *grpc.Server) {
	lis, err := net.Listen("tcp", n.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
