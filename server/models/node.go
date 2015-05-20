package models

import (
	"log"
	"net"
)

type Node struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Hostname string   `json:"hostname"`
	Uname    string   `json:"uname"`
	Addr     net.Addr `json:"addr"`
}

var ConnectedNodes []*Node

func init() {
	// Initialize Nodes Slice
	ConnectedNodes = make([]*Node, 0, 10)
}

func RegisterNode(node *Node) {
	ConnectedNodes = append(ConnectedNodes, node)
	log.Println("[Nodes] Registered connected node id:", node.Id)
}

func UnregisterNode(node *Node) {
	// Find node pos
	pos := -1
	for p, n := range ConnectedNodes {
		if n == node {
			pos = p
			break
		}
	}

	// Splice connected nodes
	if pos != -1 {
		copy(ConnectedNodes[pos:], ConnectedNodes[pos+1:])
		ConnectedNodes[len(ConnectedNodes)-1] = nil
		ConnectedNodes = ConnectedNodes[:len(ConnectedNodes)-1]
	}

	log.Println("[Nodes] Removed connected node id:", node.Id)
}
