package chord

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"sync"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
	my_list "github.com/raudel25/social-network-distributed-system/pkg/my_list"
	"google.golang.org/grpc"
)

type Node struct {
	pb.UnimplementedChordServer

	id      *big.Int
	address string

	predecessors *my_list.MyList[*Node]
	predLock     sync.RWMutex

	successors *my_list.MyList[*Node]
	sucLock    sync.RWMutex

	fingerTable FingerTable
	fingerLock  sync.RWMutex

	config *Configuration

	dictionary Storage
	dictLock   sync.RWMutex

	shutdown chan struct{}
}

func NewNode(config *Configuration, storage *Storage) *Node {
	return &Node{predecessors: my_list.NewMyList[*Node](config.SuccessorsSize), successors: my_list.NewMyList[*Node](config.SuccessorsSize),
		fingerTable: NewFingerTable(config.HashSize), dictionary: NewRamStorage(), config: config}
}

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

	if predecessor == nil {
		return nil, fmt.Errorf("not found predecessor")
	}

	return &pb.NodeResponse{
		Id:      predecessor.id.String(),
		Address: predecessor.address,
	}, nil
}

func (n *Node) GetSuccessor(ctx context.Context, req *pb.EmptyRequest) (*pb.NodeResponse, error) {
	n.sucLock.RLock()
	successor := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if successor == nil {
		return nil, fmt.Errorf("not found successor")
	}

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
		n.predecessors.SetIndex(0, newNode)
		n.predLock.Unlock()
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) Ping(ctx context.Context, req *pb.EmptyRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) Join(address string) error {
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
	n.successors.SetIndex(0, newNode)
	n.sucLock.Unlock()
	n.notify(address)

	n.predLock.Lock()
	n.predecessors.SetIndex(0, n)
	n.predLock.Unlock()

	n.fingerLock.Lock()
	n.fingerTable[0] = newNode
	n.fingerLock.Unlock()

	return nil
}

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

	return nil, fmt.Errorf("key %s\n not found", key)
}

func (n *Node) Set(ctx context.Context, req *pb.KeyValueRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.Set(req.Key, req.Value)
	n.dictLock.Unlock()

	n.sucLock.RLock()
	suc := n.successors.GetIndex(0)
	n.sucLock.RUnlock()

	if req.Rep && !equals(suc.id, n.id) {
		err := n.setReplicate(req.Key, req.Value)
		if err != nil {
			return nil, err
		}
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
		err := n.removeReplicate(req.Key)
		if err != nil {
			return nil, err
		}
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) RemoveKey(key string, value string) error {
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

func (n *Node) Start(port string) {
	// n.address = fmt.Sprintf("%s:%s", getOutboundIP().String(), port)
	n.address = fmt.Sprintf("%s:%s", "localhost", port)
	n.id = n.hashID(n.address)

	log.Printf("Starting chord server %s\n", n.address)

	s := grpc.NewServer()
	pb.RegisterChordServer(s, n)

	log.Printf("Chord server is running address:%s id:%s\n", n.address, n.id.String())

	n.createRingOrJoin()

	go n.threadListen(s)
	go n.threadStabilize()
	go n.threadCheckPredecessor()
	go n.threadCheckSuccessor()
	go n.threadFixSuccessors()
	go n.threadFixFingers()

	if port == "5002" {
		go n.threadTest()
	}
}
