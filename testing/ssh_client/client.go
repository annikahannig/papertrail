package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

func main() {
	addr := "localhost:2342"

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

	reader := bufio.NewReader(channel)
	go func() {
		text := "Hallo!\n\r"
		channel.Write([]byte(text))
		time.Sleep(10 * time.Second)
		text = "QUIT\n\r"
		channel.Write([]byte(text))
	}()

	for {
		text, err := reader.ReadString('\r')
		if err != nil {
			log.Println("[Ssh] Could not read.")
			break
		}

		text = strings.Trim(text, "\n\r")

		fmt.Println("RECV:", text, "")
	}

}
