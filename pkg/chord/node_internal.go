package chord

import (
	"crypto/sha1"
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
	} else {

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
	} else {

	}
}

func (n *Node) fixSuccessors(index int) int {
	log.Println("Fix successors")
	println(n.successors.Len())
	println(index)

	n.sucLock.RLock()
	suc := n.successors.GetIndex(index)
	len := n.successors.Len()
	last := n.successors.GetIndex(len - 1)
	n.sucLock.RUnlock()

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

	res, err := connection.client.GetSuccessor(connection.ctx, &pb.EmptyRequest{})
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

		}
	}

	return (index + 1) % n.successors.Len()
}

func (n *Node) fixFingers(index int) int {
	log.Println("Fixing finger entry")

	m := n.config.HashSize
	n.fingerLock.RLock()                         // Obtain the finger table size.
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

func (n *Node) hashID(key string) *big.Int {
	hash := sha1.New()
	hash.Write([]byte(key))
	id := new(big.Int).SetBytes(hash.Sum(nil))

	two := big.NewInt(2)
	m := big.NewInt(int64(n.config.HashSize))
	pow := big.Int{}
	pow.Exp(two, m, nil)

	id.Mod(id, &pow)
	return id
}

func (n *Node) setReplicate(key string, value string) error {
	// log.Printf("Set replicate key %s\n", key)

	// connection, err := NewGRPConnection(n.successorsFront().address)
	// if err != nil {
	// 	return err
	// }
	// defer connection.close()

	// _, err = connection.client.Set(connection.ctx, &pb.KeyValueRequest{Key: key, Value: value, Rep: false})
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (n *Node) removeReplicate(key string) error {
	// log.Printf("Remove replicate key %s\n", key)

	// connection, err := NewGRPConnection(n.successorsFront().address)
	// if err != nil {
	// 	return err
	// }
	// defer connection.close()

	// _, err = connection.client.Remove(connection.ctx, &pb.KeyRequest{Key: key, Rep: false})
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (n *Node) failPredecessorStorage(predId *big.Int) {
	// log.Println("Absorbe all predecessor data")

	// n.dictLock.RLock()
	// dict := n.dictionary.GetAll()
	// n.dictLock.RUnlock()

	// newDict := make(map[string]string)

	// n.dictLock.Lock()
	// for key, value := range dict {
	// 	keyId := n.hashID(key)
	// 	if between(keyId, predId, n.id) {
	// 		continue
	// 	}

	// 	newDict[key] = value
	// }
	// n.dictLock.Unlock()

	// connection, err := NewGRPConnection(n.successorsFront().address)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// defer connection.close()

	// connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict})
}

func (n *Node) failSuccessorStorage() {
	// log.Println("Replicate all data in new successor")

	// n.dictLock.RLock()
	// dict := n.dictionary.GetAll()
	// n.dictLock.RUnlock()

	// newDict := make(map[string]string)

	// predId := n.getPredecessorProp().id

	// for key, value := range dict {
	// 	keyId := n.hashID(key)
	// 	if !between(keyId, predId, n.id) {
	// 		continue
	// 	}

	// 	newDict[key] = value
	// }

	// connection, err := NewGRPConnection(n.successorsFront().address)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// defer connection.close()

	// connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict})
}

func (n *Node) createRing() {
	n.predLock.Lock()
	n.predecessors.SetIndex(0, n)
	n.predLock.Unlock()

	n.sucLock.Lock()
	n.successors.SetIndex(0, n)
	n.sucLock.Unlock()
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
