package chord

import (
	"context"
	"log"
	"math/big"
	"time"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
	"google.golang.org/grpc"
)

func (n *Node) findSuccessor(id *big.Int) *Node {
	findNode := n.fingerTable.FindNode(id)

	if findNode == nil {
		return n
	}

	conn, err := grpc.Dial(findNode.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChordClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.FindSuccessor(ctx, &pb.IdRequest{Id: id.String()})
	if err != nil {
		log.Fatalf("could not find successor: %v", err)
	}
	return &Node{id: hashID(res.Address), address: res.Address}
}

func (n *Node) getPredecessor(address string) (*Node, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewChordClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetPredecessor(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return &Node{id: hashID(res.Address), address: res.Address}, nil
}

func (n *Node) notify(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewChordClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client.Notify(ctx, &pb.AddressRequest{Address: n.address})
}
