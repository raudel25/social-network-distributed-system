package chord

import (
	"crypto/sha1"
	"hash"
)

type Configuration struct {
	HashFunction   func() hash.Hash // Hash function to use.
	HashSize       int              // Hash size supported
	SuccessorsSize int              // Successor to replicate and stabilize
}

// Creates and returns a default Configuration.
func DefaultConfig() *Configuration {
	HashFunction := sha1.New
	HashSize := HashFunction().Size() * 8

	config := &Configuration{
		HashFunction:   HashFunction,
		HashSize:       HashSize,
		SuccessorsSize: 5,
	}

	return config
}
