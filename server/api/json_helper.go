package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func JsonParseRequest(req *http.Request, result interface{}) error {
	// Get request body
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		log.Println("[API] Could not read request body:", err)
		return err
	}

	err = req.Body.Close()
	if err != nil {
		log.Println("[API] Error closing the request body:", err)
		return err
	}

	// Parse json
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Println("[API] Error while parsing JSON:", err)
		return err
	}

	return nil
}
