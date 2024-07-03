package chord

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (n *Node) netDiscover(port string) (string, error) {
	timeOut := 10000

	num, _ := strconv.Atoi(port)
	broadcastAddr := net.UDPAddr{
		Port: num,
		IP:   net.IPv4bcast,
	}

	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", port))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	message := []byte("Are you a chord?")
	conn.WriteTo(message, &broadcastAddr)

	buffer := make([]byte, 1024)

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Error("Error setting deadline for incoming messages.")
		return "", err
	}

	for i := 0; i < timeOut; i++ {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			continue
		}

		if string(buffer[:n]) == "I am a chord" {
			ip := strings.Split(addr.String(), ":")[0]
			log.Infof("Discover a chord in %s", ip)
			return ip, nil
		}
	}

	log.Info("Not found a chord")

	return "", nil

}
