package api

import (
	"encoding/json"
	"fmt"
	"github.com/mhannig/papertrail/server/models"
	"log"
	"net/http"
)

/**
 * GET /v1/notes
 * List all notes in DB.
 */
func NotesIndex(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /v1/notes")

	notes := models.AllNotes()
	res.Header().Set("Content-Type", "application/json")
	if len(notes) == 0 {
		res.Write([]byte("[]"))
	} else {
		json.NewEncoder(res).Encode(notes)
	}
}

/**
 * POST /v1/notes
 * Create a new note
 */
func NotesCreate(res http.ResponseWriter, req *http.Request) {
	var note models.Note

	err := JsonParseRequest(req, &note)
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
