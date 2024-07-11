package chord

import (
	"fmt"
	"log"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) getPredecessor(address string) (*Node, error) {
	connection, err := NewGRPConnection(address)
	if err != nil {
		return nil, err
	}
	defer connection.close()

	log.Printf("Find predecessor for %s", address)
	res, err := connection.client.GetPredecessor(connection.ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return &Node{id: strToBig(res.Id), address: res.Address}, nil
}

func (n *Node) notify(address string) {
	connection, err := NewGRPConnection(address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	log.Printf("Notify to %s\n", address)
	connection.client.Notify(connection.ctx, &pb.NodeRequest{Id: n.id.String(), Address: n.address})
}

func (n *Node) stabilize() {
	log.Println("Stabilizing node")

	n.sucLock.RLock()
	successor := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	pred, err := n.getPredecessor(successor.address)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if (equals(successor.id, n.id) && !equals(pred.id, n.id)) || between(pred.id, n.id, successor.id) {
		n.sucLock.Lock()
		n.successors.SetIndex(0, pred)
		n.sucLock.Unlock()
		n.notify(pred.address)
		n.replicateAllData(pred)
	} else {
		if !equals(n.id, successor.id) {
			n.notify(successor.address)
		}
	}

	log.Println("Node stabilized")
}

func (n *Node) checkSuccessor() {
	n.sucLock.RLock()
	successor := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if equals(successor.id, n.id) {
		return
	}

	log.Printf("Check successor %s\n", successor.address)

	connection, err := NewGRPConnection(successor.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	_, err = connection.client.Ping(connection.ctx, &pb.EmptyRequest{})
	if err == nil {
		return
	}

	log.Printf("Successor %s has failed\n", successor.address)

	n.sucLock.Lock()
	n.successors.RemoveIndex(0)
	n.sucLock.Unlock()

	n.sucLock.RLock()
	len := n.successors.Len()
	n.sucLock.RUnlock()

	if len == 0 {
		n.sucLock.Lock()
		n.successors.SetIndex(0, n)
		n.sucLock.Unlock()
	}
}

func (n *Node) checkPredecessor() {
	n.predLock.RLock()
	predecessor := n.predecessors.GetIndex(0)
	n.predLock.RUnlock()

	if equals(n.id, predecessor.id) {
		return
	}

	log.Printf("Check predecessor %s\n", predecessor.address)

	connection, err := NewGRPConnection(predecessor.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	_, err = connection.client.Ping(connection.ctx, &pb.EmptyRequest{})
	if err == nil {
		return
	}

	log.Printf("Predecessor %s has failed\n", predecessor.address)

	n.predLock.Lock()
	n.predecessors.RemoveIndex(0)
	n.predLock.Unlock()

	n.predLock.RLock()
	len := n.predecessors.Len()
	n.predLock.RUnlock()

	if len == 0 {
		n.predLock.Lock()
		n.predecessors.SetIndex(0, n)
		n.predLock.Unlock()
	}

	n.failPredecessorStorage(predecessor.id)
}

func (n *Node) fixSuccessors(index int) int {
	log.Println("Fix successors")

	var suc *Node
	n.sucLock.RLock()
	len := n.successors.Len()
	if index < len {
		suc = n.successors.GetIndex(index)
	}
	last := n.successors.GetIndex(len - 1)
	n.sucLock.RUnlock()

	if suc == nil {
		return 0
	}

	if suc.id == n.id && len == 1 {
		return 0
	}

	if len != 1 && equals(last.id, n.id) {
		n.sucLock.Lock()
		n.successors.RemoveIndex(len - 1)
		n.sucLock.Unlock()
		len--
	}

	connection, err := NewGRPConnection(suc.address)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	defer connection.close()

	n.sucLock.Lock()
	defer n.sucLock.Unlock()

	res, err := connection.client.GetSuccessorAndNotify(connection.ctx, &pb.NodeIndexRequest{Index: fmt.Sprintf("%d", index), Address: n.address, Id: n.id.String()})
	if err != nil {
		log.Println(err.Error())
		n.successors.RemoveIndex(index)
		if n.successors.Len() == 0 {
			n.successors.SetIndex(0, n)
		}
		return index % n.successors.Len()
	}

	sucRes := &Node{address: res.Address, id: strToBig(res.Id)}

	if equals(sucRes.id, n.id) || index == n.config.SuccessorsSize-1 {
		return 0
	}

	if index == len-1 {
		n.successors.SetIndex(index+1, sucRes)
		n.replicateAllData(sucRes)
		return (index + 1) % n.successors.Len()
	}

	sucSuc := n.successors.GetIndex(index + 1)

	if !equals(sucRes.id, sucSuc.id) {
		n.successors.SetIndex(index+1, sucRes)

		find := false

		for i := 0; i < n.successors.Len(); i++ {
			if equals(sucRes.id, n.successors.GetIndex(i).id) {
				find = true
			}
		}

		if !find {
			n.replicateAllData(sucRes)
		}
	}

	return (index + 1) % n.successors.Len()
}
