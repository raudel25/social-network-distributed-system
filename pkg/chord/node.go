package chord

import (
	"context"
	"fmt"
	"log"
	"math/big"

	// "net"
	"sync"

	"github.com/gammazero/deque"
	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
	"google.golang.org/grpc"
)

type Node struct {
	pb.UnimplementedChordServer

	id      *big.Int
	address string

	predecessor *Node
	predLock    sync.RWMutex

	successors *deque.Deque[*Node]
	sucLock    sync.RWMutex

	fingerTable FingerTable
	fingerLock  sync.RWMutex

	config *Configuration

	dictionary Storage
	dictLock   sync.RWMutex

	shutdown chan struct{}
}

func NewNode(config *Configuration, storage *Storage) *Node {
	return &Node{predecessor: nil, successors: &deque.Deque[*Node]{},
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
	predecessor := n.getPredecessorProp()

	if predecessor == nil {
		return nil, fmt.Errorf("not found predecessor")
	}

	return &pb.NodeResponse{
		Id:      predecessor.id.String(),
		Address: predecessor.address,
	}, nil
}

func (n *Node) GetSuccessor(ctx context.Context, req *pb.EmptyRequest) (*pb.NodeResponse, error) {
	successor := n.successorsFront()

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

	predecessor := n.getPredecessorProp()

	if equals(predecessor.id, n.id) || between(newNode.id, predecessor.id, n.id) {
		n.setPredecessorProp(newNode)
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

	newNode := &Node{id: strToBig(res.Id), address: res.Address}

	if err != nil {
		return err
	}

	if equals(newNode.id, n.id) {
		return fmt.Errorf("node already exists")
	}

	n.successorsPushFront(newNode)
	n.notify(address)
	n.setPredecessorProp(n)

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

	res, err := connection.client.Get(connection.ctx, &pb.KeyRequest{Key: fmt.Sprintf("key:%s", key)})
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

	_, err = connection.client.Set(connection.ctx, &pb.KeyValueRequest{Key: fmt.Sprintf("key:%s", key), Value: value})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Remove(ctx context.Context, req *pb.KeyRequest) (*pb.StatusResponse, error) {
	n.dictLock.Lock()
	n.dictionary.Remove(req.Key)
	n.dictLock.Unlock()

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

	_, err = connection.client.Remove(connection.ctx, &pb.KeyRequest{Key: fmt.Sprintf("key:%s", key)})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Start(port string) {
	// n.address = fmt.Sprintf("%s:%s", getOutboundIP().String(), port)
	n.address = fmt.Sprintf("%s:%s", "localhost", port)
	n.id = n.hashID(n.address)

	log.Printf("Starting chord server %s\n", n.address)

	s := grpc.NewServer()
	pb.RegisterChordServer(s, n)

	log.Printf("Chord server is running %s\n", n.address)

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
