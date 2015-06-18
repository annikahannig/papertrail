package sshClient

import (
	"fmt"
	"github.com/mhannig/papertrail/server/messages"
	"golang.org/x/crypto/ssh"
	"log"
)

/**
 * Protocol message handler:
 * receives from channel and modifies client shkkjjjjkktate.
 */
func handleMessages(channel ssh.Channel) error {

	// Receive and decode message from Channel
	messageId, msg, err := messages.Receive(channel)
	if err != nil {
		return err
	}

	// Implement papertrail protocol
	switch clientState {
	case S_DISCONNECTED:
		if messageId != messages.M_SERVER_INFO { // Protocol violation
			return fmt.Errorf("Received non ServerInfo message while disconnected")
		}
		serverInfo, ok := msg.(messages.MsgServerInfo)
		if ok == false {
			return fmt.Errorf("Could not decode ServerInfo Message")
		}
		// Print server info
		PrintServerInfo(serverInfo)
		PrintMotd(serverInfo)

		clientState = S_CONNECTED

		// Register / Login
		log.Println("Registering papertrail node")
		err = sendRegisterMessage(channel)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Create and send register node message
 */
func sendRegisterMessage(channel ssh.Channel) error {
	msg := messages.MsgRegisterNode{
		Name:     Node.Name,
		Hostname: Node.Uname,
		Uname:    Node.Hostname,
	}

	payload, err := messages.Encode(messages.M_REGISTER_NODE, msg)
	if err != nil {
		log.Fatal("Could not encode Node Register Failed", err)
	}

	channel.Write(payload)

	return nil
}
