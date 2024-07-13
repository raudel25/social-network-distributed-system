package chord

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) FindSuccessor(ctx context.Context, req *pb.IdRequest) (*pb.NodeResponse, error) {
	id := new(big.Int)
	id.SetString(req.Id, 10)

	successor, err := n.findSuccessor(id)
	if err != nil {
		return nil, err
	}

	return &pb.NodeResponse{
		Id:      successor.id.String(),
		Address: successor.address,
	}, nil
}

func (n *Node) GetPredecessor(ctx context.Context, req *pb.EmptyRequest) (*pb.NodeResponse, error) {
	n.predLock.RLock()
	predecessor := n.predecessors.GetIndex(0)
	n.predLock.RUnlock()

	return &pb.NodeResponse{
		Id:      predecessor.id.String(),
		Address: predecessor.address,
	}, nil
}

func (n *Node) GetSuccessorAndNotify(ctx context.Context, req *pb.NodeIndexRequest) (*pb.NodeResponse, error) {
	newNode := &Node{
		id:      strToBig(req.Id),
		address: req.Address,
	}

	n.sucLock.RLock()
	successor := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	num, err := strconv.Atoi(req.Index)
	if err != nil {
		return nil, err
	}

	n.predLock.Lock()
	if n.predecessors.Len() <= num || !equals(n.predecessors.GetIndex(num).id, newNode.id) {
		if n.predecessors.Len() < num {
			num = n.predecessors.Len()
		}
		n.predecessors.SetIndex(num, newNode)
	}
	n.predLock.Unlock()

	return &pb.NodeResponse{
		Id:      successor.id.String(),
		Address: successor.address,
	}, nil
}

func (n *Node) Notify(ctx context.Context, req *pb.NodeRequest) (*pb.StatusResponse, error) {
	newNode := &Node{
		id:      strToBig(req.Id),
		address: req.Address,
	}

	n.predLock.RLock()
	predecessor := n.predecessors.GetIndex(0)
	n.predLock.RUnlock()

	if equals(predecessor.id, n.id) || between(newNode.id, predecessor.id, n.id) {
		n.predLock.Lock()
		if equals(n.predecessors.GetIndex(0).id, n.id) {
			n.predecessors.RemoveIndex(0)
		}
		n.predecessors.SetIndex(0, newNode)
		n.predLock.Unlock()

		n.newPredecessorStorage()
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) Ping(ctx context.Context, req *pb.EmptyRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) PingLeader(ctx context.Context, req *pb.TimeRequest) (*pb.TimeResponse, error) {
	n.timeLock.Lock()
	defer n.timeLock.Unlock()

	id := req.Id
	time := req.Time

	n.time.nodeTimers[id] = time
	n.time.timeCounter = n.time.BerkleyAlgorithm()
	n.time.nodeTimers[n.id.String()] = n.time.timeCounter

	return &pb.TimeResponse{Time: n.time.timeCounter}, nil
}

func (n *Node) Election(ctx context.Context, req *pb.ElectionRequest) (*pb.NodeResponse, error) {
	selectedLeader := &Node{id: strToBig(req.SelectedLeaderId), address: req.SelectedLeaderAddress}
	firstId := strToBig(req.FirstId)

	if n.id.Cmp(selectedLeader.id) > 0 {
		selectedLeader = n
	}

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if equals(suc.id, n.id) || equals(suc.id, firstId) {
		n.leaderLock.Lock()
		n.leader = selectedLeader
		n.leaderLock.Unlock()

		return &pb.NodeResponse{Address: selectedLeader.address, Id: selectedLeader.id.String()}, nil
	}

	connection, err := NewGRPConnection(suc.address)
	if err != nil {
		return nil, err
	}
	defer connection.close()

	res, err := connection.client.Election(connection.ctx, &pb.ElectionRequest{FirstId: firstId.String(), SelectedLeaderId: selectedLeader.id.String(), SelectedLeaderAddress: selectedLeader.address})
	if err != nil {
		return nil, err
	}

	selectedLeader = &Node{id: strToBig(res.Id), address: res.Address}

	n.leaderLock.Lock()
	n.leader = selectedLeader
	n.leaderLock.Unlock()

	return &pb.NodeResponse{Address: selectedLeader.address, Id: selectedLeader.id.String()}, nil
}

func (n *Node) Join(address string, leaderAddress string) error {
	log.Printf("Joining to chord ring %s\n", address)

	connection, err := NewGRPConnection(address)
	if err != nil {
		return err
	}
	defer connection.close()

	res, err := connection.client.FindSuccessor(connection.ctx, &pb.IdRequest{Id: n.id.String()})

	if err != nil {
		return err
	}

	newNode := &Node{id: strToBig(res.Id), address: res.Address}

	if equals(newNode.id, n.id) {
		return fmt.Errorf("node already exists")
	}

	n.sucLock.Lock()
	n.successors.Clear()
	n.successors.SetIndex(0, newNode)
	n.sucLock.Unlock()

	n.predLock.Lock()
	n.predecessors.Clear()
	n.predecessors.SetIndex(0, n)
	n.predLock.Unlock()

	n.fingerLock.Lock()
	n.fingerTable[0] = newNode
	n.fingerLock.Unlock()

	n.leaderLock.Lock()
	n.leader = &Node{id: n.hashID(leaderAddress), address: leaderAddress}
	n.leaderLock.Unlock()

	n.notify(address)

	return nil
}
