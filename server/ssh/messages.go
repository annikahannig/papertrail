package sshServer

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
)

const (
	M_SERVER_INFO = iota
	M_REGISTER_NODE
	M_REGISTER_NODE_RESPONSE
	M_PRINT_JOB
	M_JOB_RESPONSE
)

type MsgServerInfo struct {
	Version          string `json:"version"`
	Motd             string `json:"motd"`
	ConnectedNodes   int    `json:"connectedNodes"`
	ConnectedClients int    `josn:"connectedClients"`
}

type MsgRegisterNode struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Uname    string `json:"uname"`
}

type MsgRegisterNodeResponse struct {
	Success bool   `json:"success"`
	Id      string `json:"id"`
}

type MsgPrintJob struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

type MsgJobResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/**
 * Encode / Decode messages
 */

func EncodeMessage(id byte, message interface{}) ([]byte, error) {
	var result []byte

	// Encode json message
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Append message type / id and length
	result = append(result, id)

	// Encode message length 4 bytes big endian
	size := uint32(len(jsonMessage))
	result = append(result,
		byte((size&0xff000000)>>24),
		byte((size&0x00ff0000)>>16),
		byte((size&0x0000ff00)>>8),
		byte((size&0x000000ff)>>0))

	// Append Json encoded message
	result = append(result, jsonMessage...)

	// Append network flush
	result = append(result, '\n', '\r')

	return result, err
}

func ReadMessage(channel ssh.Channel) (int, *interface{}) {
	return M_SERVER_INFO, nil
}
