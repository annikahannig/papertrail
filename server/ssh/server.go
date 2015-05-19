package sshServer

import (
	"golang.org/x/crypto/ssh"
	"log"
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

func (self *Server) Serve() {
	log.Println("[SSH] Starting SSH Server @", self.Listen)
}

func NewSshServer(listen string, privateKeyFilename string) *Server {
	server := Server{
		Config: &ssh.ServerConfig{
			PasswordCallback: authHandlePassword,
		},
		Listen: listen,
	}

	// Load Private Key
	privateKey := loadPrivateKey(privateKeyFile)
	server.Config.AddHostKey(privateKey)

	return nil
}
