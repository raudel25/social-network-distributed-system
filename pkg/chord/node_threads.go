package chord

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

const interval = 10

func (n *Node) threadStabilize() {
	log.Println("Stabilize thread started")

	ticker := time.NewTicker(interval * time.Second)
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

func (n *Node) threadCheckPredecessor() {
	log.Println("Check predecessor thread started")

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.checkPredecessor()
		}
	}
}

func (n *Node) threadCheckSuccessor() {
	log.Println("Check successor thread started")

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.checkSuccessor()
		}
	}
}

func (n *Node) threadFixSuccessors() {
	log.Println("Check fix successors thread started")

	next := 0 // Index of the actual successor to fix.
	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			next = n.fixSuccessors(next)
		}
	}
}

func (n *Node) threadFixFingers() {
	log.Println("Fix fingers thread started")

	next := 0                                        // Index of the actual finger entry to fix.
	ticker := time.NewTicker(interval * time.Second) // Set the time between routine activations.
	for {
		select {
		case <-n.shutdown: // If node server is shutdown, stop the thread.
			ticker.Stop()
			return
		case <-ticker.C: // If it's time, fix the correspondent finger table entry.
			next = n.fixFingers(next)
		}
	}
}

func (n *Node) threadFixStorage() {
	log.Println("Fix storage thread started")

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.fixStorage()
		}
	}
}

func (n *Node) threadListen(s *grpc.Server) {
	log.Println("Listen thread started")

	lis, err := net.Listen("tcp", n.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// BroadListen listen for broadcast messages.
func (n *Node) threadBroadListen(port string) {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Error("Error to running udp server")
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Error("Error to read the buffer")
			continue
		}

		message := string(buffer[:n])
		log.Infof("Message receive from %s: %s\n", clientAddr, message)

		if message == "Are you a chord?" {
			response := []byte("I am a chord")
			conn.WriteTo(response, clientAddr)
		}

	}

}
