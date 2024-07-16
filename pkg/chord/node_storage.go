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

func (n *Node) removeReplicate(key string, time int64) {
	log.Printf("Remove replicate key %s\n", key)

	n.sucLock.RLock()
	successors := n.successors
	defer n.sucLock.RUnlock()

	for i := 0; i < successors.Len(); i++ {
		err := n.removeReplicateNode(successors.GetIndex(i), key, time)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (n *Node) removeReplicateNode(node *Node, key string, time int64) error {
	log.Printf("Remove replicate key %s in %s\n", key, n.address)

	connection, err := NewGRPConnection(node.address)
	if err != nil {
		return err
	}
	defer connection.close()

	_, err = connection.client.Remove(connection.ctx, &pb.KeyTimeRequest{Key: key, Time: time, Rep: false})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) replicateAllData(node *Node) {
	log.Printf("Replicate all data in %s\n", node.address)

	n.dictLock.RLock()
	dict, _ := n.dictionary.GetAll()
	remove, _ := n.dictionary.GetRemoveAll()
	n.dictLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)
	newRemove := make(map[string]int64)

	n.predLock.RLock()
	pred := n.predecessors.GetIndex(0)
	n.predLock.RUnlock()

	for k, v := range dict {
		keyId := n.hashID(k)

		if between(keyId, pred.id, n.id) {
			newDict[k] = v.Value
			newVersion[k] = v.Version
		}
	}

	for key, value := range remove {
		keyId := n.hashID(key)

		if between(keyId, pred.id, n.id) {
			newRemove[key] = value.Version
		}
	}

	connection, err := NewGRPConnection(node.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict, Version: newVersion, Remove: newRemove})
}

func (n *Node) failPredecessorStorage(predId *big.Int) {
	log.Println("Absorbe all predecessor data")

	n.dictLock.RLock()
	dict, _ := n.dictionary.GetAll()
	remove, _ := n.dictionary.GetRemoveAll()
	n.dictLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)
	newRemove := make(map[string]int64)

	for key, value := range dict {
		keyId := n.hashID(key)
		if between(keyId, predId, n.id) {
			continue
		}

		newDict[key] = value.Value
		newVersion[key] = value.Version
	}

	for key, value := range remove {
		keyId := n.hashID(key)
		if between(keyId, predId, n.id) {
			continue
		}

		newRemove[key] = value.Version
	}

	n.sucLock.RLock()
	defer n.sucLock.RUnlock()

	for i := 0; i < n.successors.Len(); i++ {

		connection, err := NewGRPConnection(n.successors.GetIndex(i).address)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer connection.close()

		connection.client.SetPartition(connection.ctx, &pb.PartitionRequest{Dict: newDict, Version: newVersion, Remove: newRemove})
	}
}

func (n *Node) newPredecessorStorage() {
	log.Println("Delegate predecessor data")

	n.dictLock.RLock()
	dict, _ := n.dictionary.GetAll()
	remove, _ := n.dictionary.GetRemoveAll()
	n.dictLock.RUnlock()

	newDict := make(map[string]string)
	newVersion := make(map[string]int64)
	newRemove := make(map[string]int64)

	n.predLock.RLock()
	pred := n.predecessors.GetIndex(0)
	var predPred *Node
	if n.predecessors.Len() >= 2 {
		predPred = n.predecessors.GetIndex(1)
	} else {
		predPred = n
	}
	n.predLock.RUnlock()

	for key, value := range dict {
		keyId := n.hashID(key)
		if !between(keyId, predPred.id, pred.id) {
			continue
		}

		newDict[key] = value.Value
		newVersion[key] = value.Version
	}

	for key, value := range remove {
		keyId := n.hashID(key)
		if !between(keyId, predPred.id, pred.id) {
			continue
		}

		newRemove[key] = value.Version
	}

	connection, err := NewGRPConnection(pred.address)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer connection.close()

	res, err := connection.client.ResolveData(connection.ctx, &pb.PartitionRequest{Dict: newDict, Version: newVersion, Remove: newRemove})
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
	n.dictionary.RemoveAll(res.Remove)

}

func (n *Node) fixStorage() {
	log.Println("Fixing storage")

	n.dictLock.RLock()
	aux, _ := n.dictionary.GetAll()
	n.dictLock.RUnlock()
	log.Printf("Data storage len: %d\n", len(aux))

	n.sucLock.RLock()
	lenS := n.successors.Len()
	n.sucLock.RUnlock()

	n.predLock.Lock()
	for n.predecessors.Len() > lenS {
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

	predPred := &Node{id: strToBig(res.Id)}

	if equals(predPred.id, n.id) || equals(predPred.id, pred.id) {
		return
	}

	n.dictLock.Lock()
	dict, _ := n.dictionary.GetAll()
	defer n.dictLock.Unlock()

	for k := range dict {
		keyId := n.hashID(k)
		if between(keyId, predPred.id, n.id) {
			continue
		}

		n.timeLock.RLock()
		time := n.time.timeCounter
		n.timeLock.RUnlock()

		n.dictionary.Remove(k, time)
	}
}
