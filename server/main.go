package main

/**
 * Papertrail Server
 *
 * This server provides a RESTful HTTP API for
 * interfacing with the client (app) and provides
 * a TCP service for papertrail nodes.
 *
 * (c) 2015 Matthias Hannig
 */

import (
	"flag"
	"github.com/mhannig/papertrail/server/api"
	"github.com/mhannig/papertrail/server/config"
	"github.com/mhannig/papertrail/server/models"
	"github.com/mhannig/papertrail/server/ssh"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

/**
 * Run Papertrail and display some banner.
 */
func main() {
	log.Println("Papertrail 1.0.0                   (c) 2015 Matthias Hannig")

	// Flags
	configFilename := flag.String(
		"config",
		"./data/papertrail.conf",
		"Configuration file")

	flag.Parse()

	// Load config
	err := appconfig.Load(*configFilename)
	if err != nil {
		log.Fatal("[Config] Could not open config:", err)
	}

	// Connect to mongodb server
	session, err := mgo.Dial(appconfig.Cfg.Mongodb.Host)
	defer session.Close()

	if err != nil {
		log.Fatal("[Mongo] Could not connect to database")
	}

	session.SetMode(mgo.Monotonic, true)

	// Connect to dev database
	db := session.DB(appconfig.Cfg.Mongodb.Db)
	models.SetDatabase(db)

	// Install Timers
	go func() {
		log.Println("[Schedule] Running timed tasks")

		for {
			// Remove old sessions
			models.ClearStaleSessions()

			time.Sleep(5 * time.Minute)
		}
	}()

	// Start HTTP Server
	go func() {
		apiServer := api.NewServer(
			appconfig.Cfg.Api.Listen,
			session,
		)
		apiServer.Serve()
	}()

	// Start SSH Server
	sshServer := sshServer.NewSshServer(
		appconfig.Cfg.Ssh.Listen,
		appconfig.Cfg.Ssh.PrivateKeyFile,
	)

	sshServer.Serve()
}
