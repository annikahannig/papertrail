package messages

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

func Encode(id byte, message interface{}) ([]byte, error) {
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

	return result, err
}

func Receive(channel ssh.Channel) (uint8, interface{}, error) {
	var (
		size uint32
		id   uint8
	)

	mid := make([]byte, 1)
	msize := make([]byte, 4)

	// read first byte (message id)
	_, err := channel.Read(mid)
	if err != nil {
		return 255, nil, err
	}
	id = uint8(mid[0])

	// Read 4 bytes (message size)
	_, err = channel.Read(msize)
	if err != nil {
		return 255, nil, err
	}

	// Pack big endian size
	size = uint32(msize[0])<<24 |
		uint32(msize[1])<<16 |
		uint32(msize[2])<<8 |
		uint32(msize[3])<<0

	payload := make([]byte, size)
	_, err = channel.Read(payload)
	if err != nil {
		return 255, nil, err
	}

	// This is not elegant.
	switch id {
	case M_SERVER_INFO:
		msg := MsgServerInfo{}
		err = json.Unmarshal(payload, &msg)
		return id, msg, err
	case M_REGISTER_NODE:
		msg := MsgRegisterNode{}
		err = json.Unmarshal(payload, &msg)
		return id, msg, err
	case M_REGISTER_NODE_RESPONSE:
		msg := MsgRegisterNodeResponse{}
		err = json.Unmarshal(payload, &msg)
		return id, msg, err
	case M_PRINT_JOB:
		msg := MsgPrintJob{}
		err = json.Unmarshal(payload, &msg)
		return id, msg, err
	case M_JOB_RESPONSE:
		msg := MsgJobResponse{}
		err = json.Unmarshal(payload, &msg)
		return id, msg, err
	}
	// In fact this looks very ugly.
	// There has to be a better way. :-/

	return 255, nil, nil
}
