package chord

import (
	"math/big"
	"sync"

	"github.com/gammazero/deque"
)

type Node struct {
	id      *big.Int
	address string

	predecessor *Node
	predLock    sync.RWMutex
	successors  *deque.Deque[Node]
	sucLock     sync.RWMutex

	fingerTable FingerTable
	fingerLock  sync.RWMutex

	config *Configuration

	dictionary Storage
	dictLock   sync.RWMutex
}

func NewNode(address string, config *Configuration, storage *Storage) *Node {
	return &Node{id: hashID(address), predecessor: nil, successors: &deque.Deque[Node]{},
		fingerTable: *NewFingerTable(config.HashSize), dictionary: storage, config: config}
}
