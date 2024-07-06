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
	broadListen := flag.String("bl", "6000", "Default port broad is 6000")
	broadRequest := flag.String("br", "7000", "Default port broad is 7000")
	flag.Parse()

	config := chord.DefaultConfig()
	storage := chord.NewRamStorage()

	node := chord.NewNode(config, storage)

	node.Start(*port, *broadListen, *broadRequest)

	for {
	}

}
