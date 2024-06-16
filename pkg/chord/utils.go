package chord

import (
	"crypto/sha1"
	"math/big"
)

func hashID(key string) *big.Int {
	hash := sha1.New()
	hash.Write([]byte(key))
	return new(big.Int).SetBytes(hash.Sum(nil))
}

func between(id, start, end *big.Int) bool {
	if start.Cmp(end) < 0 {
		return id.Cmp(start) > 0 && id.Cmp(end) < 0
	} else {
		return id.Cmp(start) > 0 || id.Cmp(end) < 0
	}
}
