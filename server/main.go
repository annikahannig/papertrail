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
	"github.com/mhannig/papertrail/server/api"
	"github.com/mhannig/papertrail/server/models"
	"gopkg.in/mgo.v2"
	"log"
)

func init() {
	log.Println("Papertrail 1.0.0                   (c) 2015 Matthias Hannig")
}

func main() {

	// Connect to mongodb server
	session, err := mgo.Dial("localhost")
	defer session.Close()

	if err != nil {
		log.Fatal("[Mongo] Could not connect to database")
	}

	session.SetMode(mgo.Monotonic, true)

	// Connect to dev database
	db := session.DB("papertrail_dev")
	models.SetDatabase(db)

	server := api.NewServer(":9999", session)
	server.Serve()
}
