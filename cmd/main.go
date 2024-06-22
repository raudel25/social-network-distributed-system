package main

import (
	"flag"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
)

func main() {
	port := flag.String("p", "5000", "Default port is 5000")
	join := flag.String("j", "", "Default join is empty")
	flag.Parse()

	node := chord.NewNode(&chord.Configuration{JoinAddress: *join, SuccessorsSize: 5, HashSize: 4}, nil)

	node.Start(*port)

	for {

	}
}
