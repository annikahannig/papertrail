package main

import (
	"flag"
	"github.com/mhannig/papertrail/node/config"
	"github.com/mhannig/papertrail/node/ssh"
	"log"
	"time"
)

func main() {
	// Banner
	log.Println("Papertrail Node 1.0.0               (c) 2015 Matthias Hannig")

	// Flags
	configFilename := flag.String(
		"config",
		"./etc/papertrail-node.conf",
		"Configuration file")

	flag.Parse()

	// Load config
	err := appconfig.Load(*configFilename)
	if err != nil {
		log.Fatal("[Config] Could not open config:", err)
	}

	// Initialize node configuration
	sshClient.InitNode()

	// Main loop; reconnect if connection was interrupted
	for {
		// Connect to papertrail server
		err = sshClient.Dial(appconfig.Cfg.Server)
		if err != nil {
			log.Println(err)
		}

		log.Println("Disconnected from server: Reconnecting in 14 seconds.")
		time.Sleep(time.Second * 14)
	}
}
