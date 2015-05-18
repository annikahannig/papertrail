package api

import (
	"github.com/mhannig/papertrail/server/api/middleware"
	"github.com/mhannig/papertrail/server/models"
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
		JsonResponseError(res, 403, err, 403)
		return
	}

	JsonResponseSuccess(res, session)
}

/**
 * POST /v1/session
 */
func SessionCreate(res http.ResponseWriter, req *http.Request) {
	credentials := SessionCredentials{}
	err := JsonParseRequest(req, &credentials)
	if err != nil {
		JsonResponseError(res, 403, "invalid_credentials", 403)
		return
	}

	// Try to authenticate user
	user, err := models.AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		JsonResponseError(res, 403, err, 403)
		return
	}

	// Create new Session
	session := models.NewSession(user.Username)
	err = session.Save()
	if err != nil {
		JsonResponseError(res, 500, err, 500)
		return
	}

	JsonResponseSuccess(res, session)
}

/**
 * DELETE /v1/session
 */
func SessionDestroy(res http.ResponseWriter, req *http.Request) {
	log.Println("session delete.")
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
