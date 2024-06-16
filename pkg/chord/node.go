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
	successors  *deque.Deque[*Node]
	sucLock     sync.RWMutex

	fingerTable FingerTable
	fingerLock  sync.RWMutex

	config *Configuration

	dictionary Storage
	dictLock   sync.RWMutex

	shutdown chan struct{}
}

func NewNode(config *Configuration, storage *Storage) *Node {
	return &Node{predecessor: nil, successors: &deque.Deque[*Node]{},
		fingerTable: NewFingerTable(config.HashSize), dictionary: storage, config: config}
}

func (n *Node) FindSuccessor(ctx context.Context, req *pb.IdRequest) (*pb.AddressResponse, error) {
	id := new(big.Int)
	id.SetString(req.Id, 10)

	successor, err := n.findSuccessor(id)
	if err != nil {
		return nil, err
	}

	return &pb.AddressResponse{
		Address: successor.address,
	}, nil
}

func (n *Node) GetPredecessor(ctx context.Context, req *pb.EmptyRequest) (*pb.AddressResponse, error) {
	predecessor := n.getPredecessorProp()

	if predecessor == nil {
		return nil, fmt.Errorf("not found predecessor")
	}

	return &pb.AddressResponse{
		Address: predecessor.address,
	}, nil
}

func (n *Node) Notify(ctx context.Context, req *pb.AddressRequest) (*pb.StatusResponse, error) {
	newNode := &Node{
		id:      hashID(req.Address),
		address: req.Address,
	}

	predecessor := n.getPredecessorProp()

	if predecessor == nil || between(newNode.id, predecessor.id, n.id) {
		n.setPredecessorProp(newNode)
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) CheckPredecessor(ctx context.Context, req *pb.EmptyRequest) (*pb.StatusResponse, error) {
	predecessor := n.getPredecessorProp()

	if predecessor == nil {
		return &pb.StatusResponse{Ok: false}, nil
	}

	return &pb.StatusResponse{Ok: true}, nil
}

func (n *Node) Join(address string) error {
	log.Printf("Joining to chord ring %s\n", address)

	newNode := &Node{id: hashID(address), address: address}

	connection, err := NewGRPConnection(address)
	defer connection.close()
	if err != nil {
		return err
	}

	res, err := connection.client.FindSuccessor(connection.ctx, &pb.IdRequest{Id: hashID(address).String()})
	if err != nil {
		return err
	}

	if res.Address == n.address {
		return fmt.Errorf("node already exists")
	}

	n.successorsPushFront(newNode)
	n.notify(address)

	return nil
}

func (n *Node) Start(port string) {
	// n.address = fmt.Sprintf("%s:%s", getOutboundIP().String(), port)
	n.address = fmt.Sprintf("%s:%s", "localhost", port)
	n.id = hashID(n.address)

	log.Printf("Starting chord server %s\n", n.address)

	s := grpc.NewServer()
	pb.RegisterChordServer(s, n)

	log.Printf("Chord server is running %s\n", n.address)

	n.createRingOrJoin()

	go n.threadListen(s)
	go n.threadStabilize()
}
