package chord

import (
	"log"
	"math/big"
	"net"
)

// Checks if the given ID is between the start and end values 
//(inclusive of start, exclusive of end) in the circular ID space.
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

// Converts a string representation of a number to a *big.Int.
func strToBig(str string) *big.Int {
	bigInt := new(big.Int)

	value, _ := bigInt.SetString(str, 10)

	return value
}

// Determines the IP address that the system would use to communicate with a remote server
func getOutboundIP() net.IP {
	// Establish a UDP connection to the Google DNS server at 8.8.8.8:80.
	// This is done to determine the outbound IP address, 
	// as the connection will use the appropriate network interface to reach the remote server.
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
