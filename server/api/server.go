package api

import (
	"log"
	"mux"
	"net/http"
)

type Server struct {
	router *mux.Router
}

func NewServer(listen string) *Server {
	router := mux.NewRouter().StrictSlash(true)
	server := Server{
		router: router,
	}

	// Setup routing
	router.HandleFunc("/", welcome)

	// Create HTTP Server

}
