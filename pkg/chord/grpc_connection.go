package chord

import (
	"context"
	"log"
	"time"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
	"google.golang.org/grpc"
)

type GRPCConnection struct {
	client pb.ChordClient
	ctx    context.Context
	close  func()
}

func NewGRPConnection(address string) (*GRPCConnection, error) {
	log.Printf("Connecting to %s\n", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewChordClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &GRPCConnection{client: client, ctx: ctx, close: func() { conn.Close(); cancel() }}, nil
}
