package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type Server struct {
	Config *ssh.ServerConfig
	Listen string
}

// Helper
func authHandlePassword(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	if conn.User() == "node" && string(password) == "node" {
		return nil, nil
	}
	return nil, fmt.Errorf("[Auth] Invalid password for user:", conn.User())
}

// Load private key file (eg. id_rsa)
func loadPrivateKey(keyFile string) ssh.Signer {
	bytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal("[Ssh] Could not load private key")
	}
	key, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		log.Fatal("[Ssh] Invalid private key")
	}
	return key
}

// Handle SSH channels
func sshHandleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go func() {
			// Only allow session channeltypes
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(
					ssh.UnknownChannelType,
					fmt.Sprintf("Only accepting session channels"),
				)
				return
			}

			// conn, reqs, err := channel.Accept()
			channel, _, err := newChannel.Accept()
			if err != nil {
				log.Println("[Ssh] Accept channel failed:", err)
				return
			}

			// Send Banner
			banner := "Papertrail 1.0.0              (c) 2015 Matthias Hannig\n\r"
			channel.Write([]byte(banner))

			reader := bufio.NewReader(channel)
			for {
				text, err := reader.ReadString('\r')
				if err != nil {
					log.Println("[Ssh] Could not read.")
					break
				}

				text = strings.Trim(text, "\n\r")

				fmt.Println("RECV:", text, ";")
				if text == "QUIT" {
					break
				}

				text = text + "\n\r"

				channel.Write([]byte(text))
			}

			channel.Close()
		}()
	}
}

func NewServer(listen string, privateKeyFile string) *Server {

	server := Server{
		Config: &ssh.ServerConfig{
			PasswordCallback: authHandlePassword,
		},
		Listen: listen,
	}

	// Load Private Key
	privateKey := loadPrivateKey(privateKeyFile)
	server.Config.AddHostKey(privateKey)

	return &server
}

/**
 * Create TCP server and wait for connections
 */
func (self *Server) Start() {

	// Listen for connections
	tcpServer, err := net.Listen("tcp", self.Listen)
	if err != nil {
		log.Fatal("[Ssh] Listen on address", self.Listen, "failed.")
	}

	log.Println("[Ssh] Awaiting connections on:", self.Listen)

	// Handle connections
	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			log.Println("[Ssh] Accept failed:", err)
			continue
		}

		// Perform SSH handshake
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, self.Config)
		if err != nil {
			log.Println("[Ssh] Handshake failed:", err)
			continue
		}

		// We have a new SSH connection. Hooray.
		log.Println(
			"[Ssh] New connection from: ",
			sshConn.RemoteAddr(),
			"SSH v.",
			string(sshConn.ClientVersion()),
		)

		// Discard OOB requests
		go ssh.DiscardRequests(reqs)

		// Handle channels
		go sshHandleChannels(chans)
	}

}

/**
 * =================== [MAIN] =====================
 */
func main() {
	sshServer := NewServer(":2342", "./id_rsa")
	sshServer.Start()
}
