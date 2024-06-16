package chord

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/gammazero/deque"
	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
)

type Node struct {
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

	pb.UnimplementedChordServer
}

func NewNode(address string, config *Configuration, storage *Storage) *Node {
	return &Node{id: hashID(address), address: address, predecessor: nil, successors: &deque.Deque[*Node]{},
		fingerTable: NewFingerTable(config.HashSize), dictionary: storage, config: config}
}

func (n *Node) FindSuccessor(ctx context.Context, req *pb.IdRequest) (*pb.AddressResponse, error) {
	id := new(big.Int)
	id.SetString(req.Id, 10)
	successor := n.findSuccessor(id)
	return &pb.AddressResponse{
		Address: successor.address,
	}, nil
}

func (n *Node) GetPredecessor(ctx context.Context, req *pb.AddressRequest) (*pb.AddressResponse, error) {
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
