package chord

import (
	"context"
	"log"
	"os"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) Get(ctx context.Context, req *pb.KeyRequest) (*pb.StatusValueResponse, error) {
	n.dictLock.RLock()
	value := n.dictionary.Get(req.Key)
	n.dictLock.RUnlock()

	return &pb.StatusValueResponse{Ok: len(value) != 0, Value: value}, nil
}

func (n *Node) GetKey(key string) (*string, error) {
	log.Printf("Get key %s\n", key)

	successor, err := n.findSuccessor(n.hashID(key))
	if err != nil {
		return nil, err
	}

	connection, err := NewGRPConnection(successor.address)
	if err != nil {
		return nil, err
	}
	defer connection.close()

	res, err := connection.client.Get(connection.ctx, &pb.KeyRequest{Key: key})
	if err != nil {
		return nil, err
	}

	if res.Ok {
		return &res.Value, nil
	}

	return nil, os.ErrNotExist
}

func (n *Node) Set(ctx context.Context, req *pb.KeyValueRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.Set(req.Key, req.Value)
	n.dictLock.Unlock()

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if req.Rep && !equals(suc.id, n.id) {
		go n.setReplicate(req.Key, req.Value)
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) SetKey(key string, value string) error {
	log.Printf("Set key %s\n", key)

	successor, err := n.findSuccessor(n.hashID(key))
	if err != nil {
		return err
	}

	connection, err := NewGRPConnection(successor.address)
	if err != nil {
		return err
	}
	defer connection.close()

	_, err = connection.client.Set(connection.ctx, &pb.KeyValueRequest{Key: key, Value: value, Rep: true})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Remove(ctx context.Context, req *pb.KeyRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.Remove(req.Key)
	n.dictLock.Unlock()

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if req.Rep && !equals(n.id, suc.id) {
		go n.removeReplicate(req.Key)
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) RemoveKey(key string) error {
	log.Printf("Remove key %s\n", key)

	successor, err := n.findSuccessor(n.hashID(key))
	if err != nil {
		return err
	}

	connection, err := NewGRPConnection(successor.address)
	if err != nil {
		return err
	}
	defer connection.close()

	_, err = connection.client.Remove(connection.ctx, &pb.KeyRequest{Key: key, Rep: true})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) SetPartition(ctx context.Context, req *pb.PartitionRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.SetAll(req.Dict)
	n.dictLock.Unlock()

	return &pb.StatusResponse{Ok: true}, nil
}