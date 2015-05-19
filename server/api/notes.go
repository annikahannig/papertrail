package api

import (
	"fmt"
	"github.com/mhannig/papertrail/server/api/middleware"
	"github.com/mhannig/papertrail/server/models"
	"net/http"
)

/**
 * GET /v1/notes
 * List all notes in DB.
 */
func NotesIndex(res http.ResponseWriter, req *http.Request) {
	err := middleware.AssertAuthenticated(req)
	if err != nil {
		JsonResponseError(res, 403, err, 403)
		return
	}
	notes := models.AllNotes()
	JsonResponseSuccess(res, notes)
}

/**
 * POST /v1/notes
 * Create a new note
 */
func NotesCreate(res http.ResponseWriter, req *http.Request) {
	var note models.Note
	err := middleware.AssertAuthenticated(req)
	if err != nil {
		JsonResponseError(res, 403, err, 403)
		return
	}

	err = JsonParseRequest(req, &note)
	if err != nil {
		JsonResponseError(
			res,
			500,
			fmt.Sprintf("PARSE ERROR: %s", err),
			http.StatusNotAcceptable)
		return
	}

	err = models.InsertNote(&note)
	if err != nil {
		JsonResponseError(
			res,
			500,
			fmt.Sprintf("INSERT ERROR: %s", err),
			http.StatusNotAcceptable)
		return
	}

	JsonResponseSuccess(res, note)
}
