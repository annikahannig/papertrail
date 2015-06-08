package main

import (
	"github.com/mhannig/papertrail/server/messages"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
)

func main() {
	addr := "localhost:1992"

	user := "node"
	pass := "node"

	sshConfig := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}

	client, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		log.Fatal("[Ssh] Could not connect:", err)
	}
	defer client.Close()

	// Open channel
	channel, _, err := client.OpenChannel("session", nil)
	if err != nil {
		log.Fatal("[Ssh] Could not open channel.")
	}
	defer channel.Close()

	for {

		_, msg, err := messages.Receive(channel)
		if err != nil {
			log.Println("Error while reading message:", err)
			if err == io.EOF {
				break
			}
			continue
		}

		log.Println(*msg)

	}

	/*
		reader := bufio.NewReader(channel)
		for {
			text, err := reader.ReadString('\r')
			if err != nil {
				log.Println("[Ssh] Could not read.")
				break
			}

			text = strings.Trim(text, "\n\r")

			fmt.Println("RECV:", text, "")
		}
	*/

}
