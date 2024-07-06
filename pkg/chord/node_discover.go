package chord

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (n *Node) netDiscover(broadListen string, broadRequest string) (string, string, error) {
	log.Info("Discovering a chord ring")
	timeOut := 10000

	num, _ := strconv.Atoi(broadListen)
	broadcastAddr := net.UDPAddr{
		Port: num,
		IP:   net.IPv4bcast,
	}

	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", broadRequest))
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	message := []byte("Are you a chord?")
	conn.WriteTo(message, &broadcastAddr)

	buffer := make([]byte, 1024)

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Error("Error setting deadline for incoming messages.")
		return "", "", err
	}

	for i := 0; i < timeOut; i++ {
		nn, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			continue
		}

		res := strings.Split(string(buffer[:nn]), ";")

		if res[0] == "Yes, I am a chord" && len(res) == 2 {
			ip := strings.Split(addr.String(), ":")[0]
			log.Infof("Discover a chord in %s", ip)
			return ip, res[1], nil
		}
	}

	log.Info("Not found a chord")

	return "", "", nil

}

func (n *Node) discoverAndJoin(port string, broadListen string, broadRequest string) {
	n.leaderLock.RLock()
	leaderId := n.leader.id
	n.leaderLock.RUnlock()

	if !equals(n.id, leaderId) {
		return
	}

	discover, leaderAddress, err := n.netDiscover(broadListen, broadRequest)

	if err != nil {
		log.Error(err.Error())
		return
	}

	if discover == "" {
		return
	}

	if n.hashID(leaderAddress).Cmp(n.id) > 0 {
		n.sucLock.Lock()
		n.successors.Clear()
		n.sucLock.Unlock()

		n.predLock.Lock()
		n.predecessors.Clear()
		n.predLock.Unlock()

		err := n.Join(fmt.Sprintf("%s:%s", discover, port), leaderAddress)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
}
