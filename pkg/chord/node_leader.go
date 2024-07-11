package chord

import (
	"log"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) checkLeader() {
	n.leaderLock.RLock()
	defer n.leaderLock.RUnlock()

	if equals(n.leader.id, n.id) {
		return
	}

	log.Printf("Check leader: %s\n", n.leader.address)

	connection, err := NewGRPConnection(n.leader.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	_, err = connection.client.Ping(connection.ctx, &pb.EmptyRequest{})
	if err == nil {
		return
	}

	log.Printf("Leader has %s failed\n", n.leader.address)

	go n.electionRequested()
}

func (n *Node) electionRequested() {
	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if equals(suc.id, n.id) {
		n.leaderLock.Lock()
		n.leader = n
		n.leaderLock.Unlock()

		return
	}

	log.Println("Elections requested")

	connection, err := NewGRPConnection(suc.address)
	if err != nil {
		n.leaderLock.Lock()
		n.leader = n
		n.leaderLock.Unlock()

		log.Println(err.Error())
		return
	}
	defer connection.close()

	res, err := connection.client.Election(connection.ctx, &pb.ElectionRequest{FirstId: n.id.String(), SelectedLeaderId: n.id.String(), SelectedLeaderAddress: n.address})
	if err != nil {
		n.leaderLock.Lock()
		n.leader = n
		n.leaderLock.Unlock()

		log.Println(err.Error())
		return
	}

	n.leaderLock.Lock()
	n.leader = &Node{address: res.Address, id: strToBig(res.Id)}
	n.leaderLock.Unlock()

}
