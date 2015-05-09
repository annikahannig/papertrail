package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	listen string
	router *mux.Router
}

type ApiStats struct {
	Version string `json:"version"`
	Author  string `json:"author"`
}

/**
 * API Helper: Welcome
 * Respond with some basic API information and stats
 */
func welcome(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(ApiStats{
		Version: "1.0.0",
		Author:  "Matthias Hannig",
	})
}

/**
 * Create new HTTP API Server
 */
func NewServer(listen string) *Server {

	log.Println("Starting HTTP API Server @", listen)

	router := mux.NewRouter().StrictSlash(true)
	server := Server{
		router: router,
		listen: listen,
	}

	// Setup routing
	router.HandleFunc("/", welcome)

	return &server
}

/**
 * Create HTTP server and start serving.
 */
func (self *Server) Serve() {

	// Create HTTP Server
	err := http.ListenAndServe(self.listen, self.router)
	if err != nil {
		log.Fatal(err)
	}

}
