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
func (n *Node) threadBroadListen() {
	addr := net.UDPAddr{
		Port: 8000,
		IP:   n.ip,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error setting up UDP listener:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		message := string(buf[:n])
		if message == "DISCOVER_CHORD_NODES" {
			fmt.Printf("Received discovery message from %s\n", addr)
			conn.WriteToUDP([]byte("I am a Chord node"), addr)
		}
	}
}
