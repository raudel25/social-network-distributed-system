package chord

import "math/big"

type FingerTable []*Node

func NewFingerTable(size int) FingerTable {
	hand := make([]*Node, size)

	for i := range hand {
		hand[i] = nil
	}

	return hand
}

func (table *FingerTable) FingerId(id *big.Int, i int, m int) *big.Int {
	// Calculates 2^i.
	two := big.NewInt(2)
	pow := big.Int{}
	pow.Exp(two, big.NewInt(int64(i)), nil)

	// Calculates the sum of n and 2^i.
	sum := big.Int{}
	sum.Add(id, &pow)

	// Calculates 2^m.
	pow = big.Int{}
	pow.Exp(two, big.NewInt(int64(m)), nil)

	// Apply the mod.
	id.Mod(&sum, &pow)

	// Return the result.
	return id
}
