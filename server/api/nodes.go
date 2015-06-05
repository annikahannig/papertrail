package api

import (
	"github.com/gorilla/mux"
	"github.com/mhannig/papertrail/server/api/middleware"
	"github.com/mhannig/papertrail/server/models"
	"net/http"
)

/**
 * GET /v1/nodes
 * List all connected nodes
 */
func NodesIndex(res http.ResponseWriter, req *http.Request) {
	err := middleware.AssertAuthenticated(req)
	if err != nil {
		JsonResponseError(res, 403, err, 403)
		return
	}

	nodes := models.ConnectedNodes
	JsonResponseSuccess(res, nodes)
}

/**
 * GET /v1/nodes/:node_id
 * Show a single node
 */
func NodeShow(res http.ResponseWriter, req *http.Request) {
	err := middleware.AssertAuthenticated(req)
	if err != nil {
		JsonResponseError(res, 403, err, 403)
		return
	}

	vars := mux.Vars(req)
	nodeId := vars["nodeId"]

	if nodeId == "" {
		JsonResponseError(res, 404, "Node not found", 404)
		return
	}

	// Get a single node by id
	node := models.FindNodeById(nodeId)
	if node == nil {
		JsonResponseError(res, 404, "Node not found", 404)
		return
	}

	JsonResponseSuccess(res, node)
}
