package sshServer

import (
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

	log.Println("Received message: (", messageId, ") ", msg)

	return nil
}
