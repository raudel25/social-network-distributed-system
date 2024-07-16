package chord

import (
	"context"
	"log"
	"os"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

func (n *Node) Get(ctx context.Context, req *pb.KeyRequest) (*pb.StatusValueResponse, error) {
	n.dictLock.RLock()
	defer n.dictLock.RUnlock()

	value, err := n.dictionary.Get(req.Key)
	if err != nil {
		return nil, err
	}

	return &pb.StatusValueResponse{Ok: len(value.Value) != 0, Value: value.Value}, nil
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
	n.dictionary.Set(req.Key, Data{Value: req.Value, Version: req.Version})
	defer n.dictLock.Unlock()

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if req.Rep && !equals(suc.id, n.id) {
		n.setReplicate(req.Key, Data{Value: req.Value, Version: req.Version})
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

	n.timeLock.RLock()
	time := n.time.timeCounter
	n.timeLock.RUnlock()

	_, err = connection.client.Set(connection.ctx, &pb.KeyValueRequest{Key: key, Value: value, Version: time, Rep: true})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Remove(ctx context.Context, req *pb.KeyTimeRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.Remove(req.Key, req.Time)
	defer n.dictLock.Unlock()

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if req.Rep && !equals(n.id, suc.id) {
		n.removeReplicate(req.Key, req.Time)
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

	n.timeLock.RLock()
	time := n.time.timeCounter
	n.timeLock.RUnlock()

	_, err = connection.client.Remove(connection.ctx, &pb.KeyTimeRequest{Key: key, Time: time, Rep: true})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) SetPartition(ctx context.Context, req *pb.PartitionRequest) (*pb.StatusResponse, error) {
	newDict := make(map[string]Data)

	for key, value := range req.Dict {
		newDict[key] = Data{Value: value, Version: req.Version[key]}
	}

	n.dictLock.Lock()
	n.dictionary.SetAll(newDict)
	n.dictionary.RemoveAll(req.Remove)
	n.dictLock.Unlock()

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) ResolveData(ctx context.Context, req *pb.PartitionRequest) (*pb.ResolveDataResponse, error) {
	log.Println("Resolve data versions")
	newDict := make(map[string]Data)
	resDictValue := make(map[string]string)
	resDictVersion := make(map[string]int64)
	resRemove := make(map[string]int64)

	n.dictLock.Lock()
	defer n.dictLock.Unlock()

	dict, err := n.dictionary.GetAll()
	if err != nil {
		return nil, err
	}

	for key, value := range req.Dict {
		v, ok := dict[key]

		if ok && v.Version > req.Version[key] {
			resDictValue[key] = v.Value
			resDictVersion[key] = v.Version
		}

		if !ok || v.Version <= req.Version[key] {
			newDict[key] = Data{Value: value, Version: req.Version[key]}
		}
	}

	for key, time := range req.Remove {
		v, ok := dict[key]

		if ok && v.Version > time {
			resDictValue[key] = v.Value
			resDictVersion[key] = v.Version
		}

		if ok && v.Version <= time {
			n.dictionary.Remove(key, time)
		}
	}

	remove, _ := n.dictionary.GetRemoveAll()

	for key, value := range remove {
		t, ok := req.Version[key]

		if ok && value.Version > t {
			resRemove[key] = value.Version
		}

	}

	n.dictionary.SetAll(newDict)

	return &pb.ResolveDataResponse{Dict: resDictValue, Version: resDictVersion, Remove: resRemove}, nil
}
