package api

import (
	"encoding/json"
	"github.com/mhannig/papertrail/server/api/middleware"
	"log"
	"net/http"
)

type SessionCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/**
 * GET /v1/session
 */
func SessionShow(res http.ResponseWriter, req *http.Request) {
	session, err := middleware.CurrentSession(req)
	if err != nil {
		JsonResponseError(res, 403, "Session not authorized.", 403)
		return
	}

	json.NewEncoder(res).Encode(session)
}

/**
 * POST /v1/session
 */
func SessionCreate(res http.ResponseWriter, req *http.Request) {
	credentials := SessionCredentials{}

	log.Println("Go credentials:", credentials)
}

/**
 * DELETE /v1/session
 */
func SessionDestroy(res http.ResponseWriter, req *http.Request) {
	session, err := middleware.CurrentSession(req)
	if err != nil {
		JsonResponseError(res, 403, "Session not authorized.", 403)
		return
	}

	err = session.Destroy()
	if err != nil {
		JsonResponseError(res, 403, "session_destroy_failed", 403)
	} else {
		JsonResponseSuccess(res, "session_destroyed")
	}
}
