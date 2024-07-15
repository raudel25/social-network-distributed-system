package chord

import (
	"log"
	"math/big"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (node *Node) closestFinger(id *big.Int) *Node {
	for i := len(node.fingerTable) - 1; i >= 0; i-- {
		node.fingerLock.RLock()
		finger := node.fingerTable[i]
		node.fingerLock.RUnlock()

		if finger != nil && between(finger.id, node.id, id) {
			return finger
		}
	}

	return node
}

func (n *Node) findSuccessor(id *big.Int) (*Node, error) {
	log.Printf("Find successor for %s", id.String())

	findNode := n.closestFinger(id)

	if equals(findNode.id, n.id) {
		n.sucLock.RLock()
		defer n.sucLock.RUnlock()
		return n.predecessors.GetIndex(0), nil
	}

	connection, err := NewGRPConnection(findNode.address)
	if err != nil {
		return nil, err
	}
	defer connection.close()

	res, err := connection.client.FindSuccessor(connection.ctx, &pb.IdRequest{Id: id.String()})
	if err != nil {
		return nil, err
	}

	return &Node{id: strToBig(res.Id), address: res.Address}, nil
}

func (n *Node) fixFingers(index int) int {
	log.Println("Fixing finger entry")

	m := n.config.HashSize // Obtain the finger table size.
	n.fingerLock.RLock()
	id := n.fingerTable.FingerId(n.id, index, m) // Obtain node.ID + 2^(next) mod(2^m).
	n.fingerLock.RUnlock()
	suc, err := n.findSuccessor(id) // Obtain the node that succeeds ID = node.ID + 2^(next) mod(2^m).

	// In case of error finding the successor, report the error and skip this finger.
	if err != nil || suc == nil {
		log.Printf("Successor of ID not found.This finger fix was skipped %s\n", err.Error())

		return (index + 1) % m
	}

	log.Printf("Correspondent finger found at %s\n", suc.address)

	// If the successor of this ID is this node, then the ring has already been turned around.
	// Clean the remaining positions and return index 0 to restart the fixing cycle.
	if equals(suc.id, n.id) {
		for i := index; i < m; i++ {
			n.fingerLock.Lock()    // Lock finger table to write on it, and unlock it after.
			n.fingerTable[i] = nil // Clean the correspondent position on the finger table.
			n.fingerLock.Unlock()
		}
		return 0
	}

	n.fingerLock.Lock()
	n.fingerTable[index] = suc // Update the correspondent position on the finger table.
	n.fingerLock.Unlock()

	// Return the next index to fix.
	return (index + 1) % m
}
