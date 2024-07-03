package main

import (
	"flag"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
	"github.com/raudel25/social-network-distributed-system/pkg/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SettingLogger(log.DebugLevel, ".")
	port := flag.String("p", "5000", "Default port is 5000")
	broad := flag.String("b", "8000", "Default port broad is 8000")
	flag.Parse()

	node := chord.NewNode(&chord.Configuration{SuccessorsSize: 5, HashSize: 4}, nil)

	node.Start(*port, *broad)

	for {
	}

}
