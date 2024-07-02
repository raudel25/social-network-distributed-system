package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
	"github.com/raudel25/social-network-distributed-system/pkg/logging"
)

func main() {
	logging.SettingLogger(log.DebugLevel, ".")
	port := flag.String("p", "5000", "Default port is 5000")
	join := flag.String("j", "", "Default join is empty")
	flag.Parse()

	node := chord.NewNode(&chord.Configuration{JoinAddress: *join, SuccessorsSize: 5, HashSize: 4}, nil)

	node.Start(*port)

	for {
	}
}
