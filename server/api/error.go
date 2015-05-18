package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(res http.ResponseWriter, code int, message string, status int) {
	msg, err := json.Marshal(ErrorResponse{
		Code:    code,
		Message: message})
	if err != nil {
		log.Println("[Error] Could not encode error:", err)
		return
	}
	http.Error(res, string(msg), status)
}
