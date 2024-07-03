package chord

import (
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"
)

func (n *Node) hashID(key string) *big.Int {
	hash := n.config.HashFunction()
	hash.Write([]byte(key))
	id := new(big.Int).SetBytes(hash.Sum(nil))

	two := big.NewInt(2)
	m := big.NewInt(int64(n.config.HashSize))
	pow := big.Int{}
	pow.Exp(two, m, nil)

	id.Mod(id, &pow)
	return id
}

func (n *Node) createRing() {
	log.Info("Create a chord ring")
	n.predLock.Lock()
	n.predecessors.SetIndex(0, n)
	n.predLock.Unlock()

	n.sucLock.Lock()
	n.successors.SetIndex(0, n)
	n.sucLock.Unlock()
}

func (n *Node) createRingOrJoin(broad string, port string) {
	discover, err := n.netDiscover(broad)

	if err == nil && discover != "" {
		err := n.Join(fmt.Sprintf("%s:%s", discover, port))
		if err != nil {
			log.Info(err.Error())
		}
		return
	}

	n.createRing()
}
