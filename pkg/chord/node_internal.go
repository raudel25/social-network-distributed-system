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
	if err != nil {
		return nil, err
	}
	defer connection.close()

	res, err := connection.client.FindSuccessor(connection.ctx, &pb.IdRequest{Id: id.String()})
	if err != nil {
		return nil, err
	}

	return &Node{id: hashID(res.Address), address: res.Address}, nil
}

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

	return &Node{id: hashID(res.Address), address: res.Address}, nil
}

func (n *Node) notify(address string) {
	connection, err := NewGRPConnection(address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	log.Printf("Notify to %s\n", address)
	connection.client.Notify(connection.ctx, &pb.AddressRequest{Address: n.address})
}

func (n *Node) stabilize() {
	log.Println("Stabilizing node")

	successor := n.successorsFront()
	pred, err := n.getPredecessor(successor.address)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if between(pred.id, n.id, successor.id) {
		n.successorsPushFront(pred)
		n.notify(pred.address)
	} else {
		if n.address != successor.address {
			n.notify(successor.address)
		}
	}

	log.Println("Node stabilized")
}

func (n *Node) checkSuccessor() {
	successor := n.successorsFront()

	if successor.address == n.address {
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

	n.successorsPopFront()
	if n.successorsLen() == 0 {
		n.successors.PushBack(n)
	}
}

func (n *Node) checkPredecessor() {
	predecessor := n.getPredecessorProp()

	if predecessor.address == n.address {
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
	n.setPredecessorProp(n)
}

func (n *Node) createRing() {
	n.successorsPushBack(n)
	n.setPredecessorProp(n)
}

func (n *Node) createRingOrJoin() {
	if n.config.JoinAddress != "" {
		err := n.Join(n.config.JoinAddress)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	n.createRing()
}
