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
