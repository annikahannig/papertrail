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
	"github.com/BurntSushi/toml"
	"github.com/mhannig/papertrail/server/api"
	"github.com/mhannig/papertrail/server/models"
	"gopkg.in/mgo.v2"
	"log"
)

/**
 * Config
 */
type config struct {
	Listen  string
	Mongodb mongodbConfig
}

type mongodbConfig struct {
	Host string
	Db   string
}

/**
 * Initialize Papertrail and display some banner.
 */
func init() {
	log.Println("Papertrail 1.0.0               (c) 2015 Matthias Hannig")
}

func main() {

	// Flags
	configFilename := flag.String(
		"config",
		"./data/papertrail.conf",
		"Configuration file")

	flag.Parse()

	// Load config
	var cfg config
	_, err := toml.DecodeFile(*configFilename, &cfg)
	if err != nil {
		log.Fatal("[Config] Could not open config:", err)
	}

	// Connect to mongodb server
	session, err := mgo.Dial(cfg.Mongodb.Host)
	defer session.Close()

	if err != nil {
		log.Fatal("[Mongo] Could not connect to database")
	}

	session.SetMode(mgo.Monotonic, true)

	// Connect to dev database
	db := session.DB(cfg.Mongodb.Db)
	models.SetDatabase(db)

	server := api.NewServer(cfg.Listen, session)
	server.Serve()
}
