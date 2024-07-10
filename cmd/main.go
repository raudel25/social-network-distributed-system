package main

import (
	"flag"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
	"github.com/raudel25/social-network-distributed-system/pkg/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SettingLogger(log.DebugLevel, ".")
	port := flag.String("p", "10000", "Default port is 10000")
	broadListen := flag.String("bl", "11000", "Default port broad is 11000")
	broadRequest := flag.String("br", "12000", "Default port broad is 12000")
	flag.Parse()

	config := chord.DefaultConfig()
	storage := chord.NewRamStorage()

	node := chord.NewNode(config, storage)

	node.Start(*port, *broadListen, *broadRequest)

	for {
	}

}
