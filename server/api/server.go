package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mhannig/papertrail/server/api/middleware"
	"github.com/mhannig/papertrail/server/config"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type Server struct {
	listen  string
	router  *mux.Router
	session *mgo.Session
}

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
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
func NewServer(listen string, mongoSession *mgo.Session) *Server {

	log.Println("Starting HTTP API Server @", listen)

	router := mux.NewRouter().StrictSlash(true)
	server := Server{
		router:  router,
		listen:  listen,
		session: mongoSession,
	}

	// Setup routing
	routes := []Route{
		Route{
			"GET", "/", welcome,
		},
		Route{
			"GET", "/v1", welcome,
		},
		// Sessions
		Route{
			"GET", "/v1/session", SessionShow,
		},
		Route{
			"POST", "/v1/session", SessionCreate,
		},
		Route{
			"DELETE", "/v1/session", SessionDestroy,
		},
		// Notes
		Route{
			"GET", "/v1/notes", NotesIndex,
		},
		Route{
			"POST", "/v1/notes", NotesCreate,
		},
		// Connected Nodes
		Route{
			"GET", "/v1/nodes", NodesIndex,
		},
		Route{
			"GET", "/v1/nodes/{nodeId}", NodeShow,
		},
	}

	debugRoutes := []Route{
		Route{
			"GET", "/v1/sessions", SessionsIndex,
		},
	}

	log.Println("Installing routes:")
	for _, route := range routes {
		log.Println(route.Method, "\t", route.Path)
		router.
			Methods(route.Method).
			Path(route.Path).
			Handler(middleware.AuthSessionToken(route.Handler))
	}

	if appconfig.Cfg.Debug {
		log.Println("Installing debug routes:")
		for _, route := range debugRoutes {
			log.Println(route.Method, "\t", route.Path)
			router.
				Methods(route.Method).
				Path(route.Path).
				Handler(middleware.AuthSessionToken(route.Handler))
		}
	}

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
