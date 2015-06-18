package sshClient

import (
	"github.com/mhannig/papertrail/node/config"
	"log"
	"os"
	"os/exec"
)

type PapertrailNode struct {
	Name     string
	Hostname string
	Uname    string
}

func NewPapertrailNode() *PapertrailNode {

	// Get hostname and uname
	hostname, _ := os.Hostname()
	uname, err := exec.Command("uname", "-a").Output()
	if err != nil {
		log.Fatal("Could not get uname:", err)
	}

	// Create node struct
	node := PapertrailNode{
		Name:     appconfig.Cfg.NodeName,
		Hostname: hostname,
		Uname:    string(uname),
	}

	return &node
}

// Global Node
var Node *PapertrailNode

func InitNode() {
	Node = NewPapertrailNode()
}
