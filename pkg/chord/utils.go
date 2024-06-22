package chord

import (
	"crypto/sha1"
	"log"
	"math/big"
	"net"
)

func hashID(key string) *big.Int {
	hash := sha1.New()
	hash.Write([]byte(key))
	return new(big.Int).SetBytes(hash.Sum(nil))
}

func between(id, start, end *big.Int) bool {
	if start.Cmp(end) <= 0 {
		return id.Cmp(start) > 0 && id.Cmp(end) < 0
	} else {
		return id.Cmp(start) > 0 || id.Cmp(end) < 0
	}
}

func equals(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) == 0
}

func strToBig(str string) *big.Int {
	bigInt := new(big.Int)

	value, _ := bigInt.SetString(str, 10)

	return value
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
