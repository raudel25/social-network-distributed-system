package chord

import (
	"log"
	"math/big"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) findSuccessor(id *big.Int) (*Node, error) {
	log.Printf("Find successor for %s", id.String())

	n.fingerLock.Lock()
	findNode := n.fingerTable.FindNode(id)
	n.fingerLock.Unlock()

	if findNode == nil {
		return n, nil
	}

	connection, err := NewGRPConnection(n.address)
	defer connection.close()
	if err != nil {
		return nil, err
	}

	res, err := connection.client.FindSuccessor(connection.ctx, &pb.IdRequest{Id: id.String()})
	if err != nil {
		return nil, err
	}

	return &Node{id: hashID(res.Address), address: res.Address}, nil
}

func (n *Node) getPredecessor(address string) (*Node, error) {
	connection, err := NewGRPConnection(address)
	defer connection.close()
	if err != nil {
		return nil, err
	}

	log.Printf("Find predecessor for %s", address)
	res, err := connection.client.GetPredecessor(connection.ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return &Node{id: hashID(res.Address), address: res.Address}, nil
}

func (n *Node) notify(address string) error {
	connection, err := NewGRPConnection(address)
	defer connection.close()
	if err != nil {
		return err
	}

	log.Printf("Notify to %s\n", address)
	connection.client.Notify(connection.ctx, &pb.AddressRequest{Address: n.address})

	return nil
}

func (n *Node) stabilize() {
	log.Println("Stabilizing node")

	successor := n.successorsFront()
	pred, err := n.getPredecessor(successor.address)

	if err == nil && between(pred.id, n.id, successor.id) {
		n.successorsPushFront(pred)
	}

	n.notify(n.successorsFront().address)

	log.Println("Node stabilized")
}
