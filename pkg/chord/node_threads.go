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
func (n *Node) threadBroadListen(broad string) {
	// Wait for the specific port to be free to use.
	pc, err := net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%s", broad))
	for err != nil {
		pc, err = net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%s", broad))
	}
	// Close the listening socket at the end of function.
	defer func(pc net.PacketConn) {
		err := pc.Close()
		if err != nil {
			return
		}
	}(pc)

	// Start listening messages.
	for {
		// If node server is shutdown, return.
		if !isOpen(n.shutdown) {
			return
		}

		// Create the buffer to store the message.
		buf := make([]byte, 1024)
		// Wait for a message.
		n, address, err := pc.ReadFrom(buf)
		if err != nil {
			log.Errorf("Incoming broadcast message error.\n%s", err.Error())
			continue
		}

		log.Debugf("Incoming response message. %s sent this: %s", address, buf[:n])

		// If the incoming message is the specified one, answer with the specific response.
		if string(buf[:n]) == "Chord?" {
			_, err = pc.WriteTo([]byte("I am chord"), address)
			if err != nil {
				log.Errorf("Error responding broadcast message.\n%s", err.Error())
				continue
			}
		}
	}
}
