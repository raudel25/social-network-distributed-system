package chord

import "math/big"

type FingerTable []*Node

func NewFingerTable(size int) FingerTable {
	hand := make([]*Node, size) // Build the new array of fingers.

	for i := range hand {
		hand[i] = nil
	}

	return hand
}

func (table FingerTable) FindNode(id *big.Int) *Node {
	var node *Node

	for _, n := range table {
		if id.Cmp(n.id) != 1 {
			node = n
		}
	}

	return node
}
