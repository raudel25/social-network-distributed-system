package chord

import (
	"fmt"
	"log"
	"math/big"
	"net"

	"sync"

	pb "github.com/raudel25/social-network-distributed-system/pkg/chord/grpc"
	my_list "github.com/raudel25/social-network-distributed-system/pkg/my_list"
	"google.golang.org/grpc"
)

type Node struct {
	pb.UnimplementedChordServer

	id      *big.Int
	address string
	ip      net.IP

	time     *NodeTime
	timeLock sync.RWMutex

	leader     *Node
	leaderLock sync.RWMutex

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

func NewNode(config *Configuration, storage Storage) *Node {
	return &Node{
		predecessors: my_list.NewMyList[*Node](config.SuccessorsSize),
		successors:   my_list.NewMyList[*Node](config.SuccessorsSize),
		fingerTable:  NewFingerTable(config.HashSize),
		dictionary:   storage,
		config:       config,
	}
}

// DefaultNode creates and returns a new Node with default configurations.
func DefaultNode() (*Node, error) {
	conf := DefaultConfig()                                      // Creates a default configuration.
	dictionary := NewDictStorage("/tmp/socialnetwork-data.json") // Creates a default dictionary.
	return NewNode(conf, dictionary), nil
}

func (n *Node) Start(port string, broadListen string, broadRequest string) {
	n.ip = getOutboundIP()
	n.address = fmt.Sprintf("%s:%s", n.ip.String(), port)

	n.id = n.hashID(n.address)
	n.time = NewNodeTime(n.id)

	log.Printf("Starting chord server %s\n", n.address)

	s := grpc.NewServer()
	pb.RegisterChordServer(s, n)

	log.Printf("Chord server is running address:%s id:%s\n", n.address, n.id.String())

	go n.threadListen(s)

	n.createRingOrJoin(broadListen, broadRequest, port)

	go n.threadStabilize()
	go n.threadCheckPredecessor()
	go n.threadCheckSuccessor()
	go n.threadFixSuccessors()
	go n.threadFixFingers()
	go n.threadFixStorage()
	go n.threadDiscoverAndJoin(port, broadListen, broadRequest)
	go n.threadCheckLeader()
	go n.threadBroadListen(broadListen)
	go n.threadRequestElections()
	go n.threadUpdateTime()
}
