package chord

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

const interval = 10
const intervalS = 60

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

func (n *Node) threadCheckLeader() {
	log.Println("Check leader thread started")

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.checkLeader()
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

func (n *Node) threadDiscoverAndJoin(port string, broadListen string, broadRequest string) {
	log.Println("Fix storage thread discover and join started")

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.discoverAndJoin(port, broadListen, broadRequest)
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
		nn, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Error("Error to read the buffer")
			continue
		}

		n.leaderLock.RLock()
		leaderId := n.leader.id
		n.leaderLock.RUnlock()

		if !equals(leaderId, n.id) {
			continue
		}

		message := string(buffer[:nn])
		log.Infof("Message receive from %s: %s", clientAddr, message)

		if message == "Are you a chord?" {
			n.leaderLock.RLock()
			leader := n.leader
			n.leaderLock.RUnlock()

			response := []byte(fmt.Sprintf("Yes, I am a chord;%s", leader.address))
			conn.WriteTo(response, clientAddr)
		}

	}
}

func (n *Node) threadRequestElections() {
	log.Println("Election requested thread started")

	ticker := time.NewTicker(intervalS * time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.leaderLock.RLock()
			leaderId := n.leader.id
			n.leaderLock.RUnlock()

			if equals(n.id, leaderId) {
				n.electionRequested()
			}
		}
	}

}

func (n *Node) threadUpdateTime() {
	log.Println("Update time thread started")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-n.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			n.timeLock.Lock()
			n.time.timeCounter += 1
			n.time.nodeTimers[n.id.String()] += 1
			n.timeLock.Unlock()
		}
	}

}
