package chord

import (
	"math/big"
	"time"
)

type NodeTime struct {
	timeCounter int64
	nodeTimers  map[string]int64
}

func NewNodeTime(id *big.Int) *NodeTime {
	now := time.Now().Unix()
	m := make(map[string]int64)
	m[id.String()] = now

	return &NodeTime{timeCounter: now, nodeTimers: m}
}

func (nt *NodeTime) BerkleyAlgorithm() int64 {
	var w int64
	w = 0

	for _, v := range nt.nodeTimers {
		w += v
	}

	return w / int64(len(nt.nodeTimers))
}
