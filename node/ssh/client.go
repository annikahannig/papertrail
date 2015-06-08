package sshClient

import (
	"fmt"
	"github.com/mhannig/papertrail/server/messages"
	"golang.org/x/crypto/ssh"
	"io"
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

var state uint8

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

		// Receive and decode message from Channel
		messageId, msg, err := messages.Receive(channel)
		if err != nil {
			log.Println("Error while reading message:", err)
			if err == io.EOF {
				break
			}
			continue
		}

		// Implement papertrail protocol
		switch state {
		case S_DISCONNECTED:
			if messageId != messages.M_SERVER_INFO { // Protocol violation
				return fmt.Errorf("Received non ServerInfo message while disconnected")
			}
			serverInfo, ok := msg.(messages.MsgServerInfo)
			if ok == false {
				continue
			}
			// Print server info
			PrintServerInfo(serverInfo)
			PrintMotd(serverInfo)
		}

	}

	return nil
}
