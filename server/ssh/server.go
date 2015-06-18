package sshServer

import (
	"fmt"
	"github.com/mhannig/papertrail/server/messages"
	"github.com/mhannig/papertrail/server/models"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

type Server struct {
	Config *ssh.ServerConfig
	Listen string
}

/*
 * == HELPER
 */

// Handle password auath
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

/**
 * Create new SSH Server
 */
func NewSshServer(listen string, privateKeyFilename string) *Server {
	server := Server{
		Config: &ssh.ServerConfig{
			PasswordCallback: authHandlePassword,
		},
		Listen: listen,
	}

	// Load Private Key
	privateKey := loadPrivateKey(privateKeyFilename)
	server.Config.AddHostKey(privateKey)

	return &server
}

/**
 * SSH Server Main
 */
func (self *Server) Serve() {
	log.Println("[SSH] Starting SSH Server @", self.Listen)
	tcpServer, err := net.Listen("tcp", self.Listen)
	if err != nil {
		log.Fatal("[SSH] Could not listen on configured address:", err)
	}

	// Handle TCP connections
	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			log.Println("[SSH] Accept() failed:", err)
			continue // accept next connection
		}

		// Perform SSH handshake
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, self.Config)
		if err != nil {
			log.Println(
				"[SSH] Handshake failed for connection:",
				conn.RemoteAddr(), " - ", err,
			)
			continue
		}

		// Everything is fine: We have a TCP connection and performed
		// a succesfull SSH handshake.
		log.Println(
			"[SSH] New connection from:", sshConn.RemoteAddr(),
			"SSH v.", string(sshConn.ClientVersion()),
		)

		// Connect node
		node := &models.Node{
			Addr:    sshConn.RemoteAddr(),
			SshConn: sshConn,
		}

		// Discard OOB requests
		go ssh.DiscardRequests(reqs)

		// Handle SSH channels
		go self.handleChannels(node, chans)
	}
}

/**
 * Establish SSH channels
 */
func (self *Server) handleChannels(node *models.Node, channels <-chan ssh.NewChannel) {
	for newChannel := range channels {
		// Only allow session channeltypes
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(
				ssh.UnknownChannelType,
				"ERROR: The server only accepts session channels",
			)
			continue
		}

		// Accept this channel
		channel, _, err := newChannel.Accept()
		if err != nil {
			log.Println("[SSH] Channel Accept() failed:", err)
			continue
		}

		// Everything is fine: We have a TCP connection,
		// a successfully established SSH connection
		// and an open SSH channel.
		go self.handleSession(node, channel)
	}
}

/**
 * Handle a single ssh session on a channel
 */
func (self *Server) handleSession(node *models.Node, channel ssh.Channel) {

	// BYTES[0]:    message ID   (8bit)
	// BYTES[1..4]: message Size (32bit)
	// BYTES[5..]:  packed (protobuf) message (4294967296 bytes max)

	// Send server info
	i := 0
	info := messages.MsgServerInfo{
		"1.0.2",
		"Papertrail for the win.",
		i,
		1,
	}

	msg, _ := messages.Encode(messages.M_SERVER_INFO, info)
	channel.Write(msg)

	for {
		err := handleMessages(channel)
		if err != nil {
			log.Println("Could not read message:", err)
			break
		}
	}
}
