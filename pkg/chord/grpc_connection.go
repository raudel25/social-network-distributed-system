package chord

import (
	"context"
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewChordClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	return &GRPCConnection{client: client, ctx: ctx, close: func() { conn.Close(); cancel() }}, nil
}
