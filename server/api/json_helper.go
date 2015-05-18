package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonErrorResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

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

func JsonResponseSuccess(res http.ResponseWriter, message interface{}) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(message)
}

func JsonResponseError(res http.ResponseWriter, code int, message interface{}, status int) {
	res.Header().Set("Content-Type", "application/json")
	msg, err := json.Marshal(JsonErrorResponse{
		Code:    code,
		Message: fmt.Sprintf("%s", message)})
	if err != nil {
		log.Println("[Error] Could not encode error:", err)
		return
	}
	http.Error(res, string(msg), status)
}
