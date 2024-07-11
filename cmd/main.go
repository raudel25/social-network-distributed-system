package main

import (
	socialnetwork "github.com/raudel25/social-network-distributed-system/internal/services"
	"github.com/raudel25/social-network-distributed-system/pkg/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SettingLogger(log.DebugLevel, ".")
	rsaPrivateKeyPath := "pv.pem"
	rsaPublicteKeyPath := "pub.pem"
	network := "tcp"

	socialnetwork.Start(rsaPrivateKeyPath, rsaPublicteKeyPath, network)
	
	for {
	}

}
