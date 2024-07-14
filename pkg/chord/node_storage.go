package chord

import (
	"log"
	"math/big"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) setReplicate(key string, value Data) {
	log.Printf("Set replicate key %s\n", key)

	n.sucLock.RLock()
	successors := n.successors
	defer n.sucLock.RUnlock()

	for i := 0; i < successors.Len(); i++ {
		err := n.setReplicateNode(successors.GetIndex(i), key, value)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (n *Node) setReplicateNode(node *Node, key string, value Data) error {
	log.Printf("Set replicate key %s in %s\n", key, node.address)

	connection, err := NewGRPConnection(node.address)
	if err != nil {
		return err
	}
	defer connection.close()

	_, err = connection.client.Set(connection.ctx, &pb.KeyValueRequest{Key: key, Value: value.Value, Version: value.Version, Rep: false})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) removeReplicate(key string) {
	log.Printf("Remove replicate key %s\n", key)

	n.sucLock.RLock()
	successors := n.successors
	defer n.sucLock.RUnlock()

	for i := 0; i < successors.Len(); i++ {
		err := n.removeReplicateNode(successors.GetIndex(i), key)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (n *Node) removeReplicateNode(node *Node, key string) error {
	log.Printf("Remove replicate key %s in %s\n", key, n.address)

	connection, err := NewGRPConnection(node.address)
	if err != nil {
		return err
	}
	defer connection.close()

	_, err = connection.client.Remove(connection.ctx, &pb.KeyRequest{Key: key, Rep: false})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) replicateAllData(node *Node) {
	log.Printf("Replicate all data in %s\n", node.address)

	n.dictLock.Lock()
	dict := n.dictionary.GetAll()
	defer n.dictLock.Unlock()

	n.predLock.RLock()
	pred := n.predecessors.GetIndex(0)
	n.predLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)

	for k, v := range dict {
		keyId := n.hashID(k)

		if between(keyId, pred.id, n.id) {
			newDict[k] = v.Value
			newVersion[k] = v.Version
		}
	}

	connection, err := NewGRPConnection(node.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict, Version: newVersion})
}

func (n *Node) failPredecessorStorage(predId *big.Int) {
	log.Println("Absorbe all predecessor data")

	n.dictLock.RLock()
	dict := n.dictionary.GetAll()
	n.dictLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)

	n.dictLock.Lock()
	defer n.dictLock.Unlock()

	for key, value := range dict {
		keyId := n.hashID(key)
		if between(keyId, predId, n.id) {
			continue
		}

		newDict[key] = value.Value
		newVersion[key] = value.Version
	}

	n.sucLock.RLock()
	defer n.sucLock.RUnlock()

	connection, err := NewGRPConnection(n.successors.GetIndex(n.successors.Len() - 1).address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict})
}

func (n *Node) newPredecessorStorage() {
	log.Println("Delegate predecessor data")

	n.dictLock.RLock()
	dict := n.dictionary.GetAll()
	n.dictLock.RUnlock()

	n.predLock.RLock()
	pred := n.predecessors.GetIndex(0)
	var predPred *Node
	if n.predecessors.Len() >= 2 {
		predPred = n.predecessors.GetIndex(1)
	} else {
		predPred = n
	}
	n.predLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)

	for key, value := range dict {
		keyId := n.hashID(key)
		if !between(keyId, predPred.id, pred.id) {
			continue
		}

		newDict[key] = value.Value
		newVersion[key] = value.Version
	}

	connection, err := NewGRPConnection(pred.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	res, err := connection.client.ResolveData(connection.ctx, &pb.PartitionRequest{Dict: newDict})
	if err != nil {
		log.Println(err.Error())
		return
	}

	newResDict := make(map[string]Data)

	for k, v := range res.Dict {
		newResDict[k] = Data{Value: v, Version: res.Version[k]}
	}

	n.dictLock.Lock()
	defer n.dictLock.Unlock()
	n.dictionary.SetAll(newResDict)

}

func (n *Node) fixStorage() {
	log.Println("Fixing storage")

	n.sucLock.RLock()
	len := n.successors.Len()
	n.sucLock.RUnlock()

	n.predLock.Lock()
	for n.predecessors.Len() > len {
		n.predecessors.RemoveIndex(n.predecessors.Len() - 1)
		if n.predecessors.Len() == 0 {
			n.predecessors.SetIndex(0, n)
			break
		}
	}
	n.predLock.Unlock()

	n.predLock.RLock()
	pred := n.predecessors.GetIndex(n.predecessors.Len() - 1)
	n.predLock.RUnlock()

	if equals(pred.id, n.id) {
		return
	}

	connection, err := NewGRPConnection(pred.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	res, err := connection.client.GetPredecessor(connection.ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	pred = &Node{id: strToBig(res.Id)}

	if equals(pred.id, n.id) {
		return
	}

	n.dictLock.Lock()
	dict := n.dictionary.GetAll()
	defer n.dictLock.Unlock()

	for k := range dict {
		keyId := n.hashID(k)
		if between(keyId, pred.id, n.id) {
			continue
		}

		n.dictionary.Remove(k)
	}
}
