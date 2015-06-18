package sshClient

import (
	"github.com/mhannig/papertrail/server/messages"
	"golang.org/x/crypto/ssh"
	"log"
)

/**
 * Client states
 */
const (
	S_DISCONNECTED = iota
	S_REGISTERED
	S_CONNECTED
	S_READY
)

var clientState uint8

// Helper: Print Server Info
func PrintServerInfo(info messages.MsgServerInfo) {
	log.Println(
		"Connected to Papertrail Server v.", info.Version,
		"as one of", info.ConnectedNodes, "connected nodes.",
	)
	log.Println(
		"There are", info.ConnectedClients, "clients connected.",
	)
}

func PrintMotd(info messages.MsgServerInfo) {
	log.Println("---------------------[MOTD]-----------------------")
	log.Println(info.Motd)
	log.Println("--------------------------------------------------")
}

func Dial(addr string) error {
	user := "node"
	pass := "node"

	// Start papertrail node
	clientState = S_DISCONNECTED

	sshConfig := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}

	client, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	// Open channel
	channel, _, err := client.OpenChannel("session", nil)
	if err != nil {
		return err
	}
	defer channel.Close()

	// Main receive and respond loop
	for {
		err := handleMessages(channel)
		if err != nil {
			log.Println(err)
			break // Close connection
		}
	}

	return nil
}
