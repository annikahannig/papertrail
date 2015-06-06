package sshServer

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
