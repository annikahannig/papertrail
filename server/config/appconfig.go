package appconfig

import (
	"github.com/BurntSushi/toml"
)

/**
 * Config
 */
type AppConfig struct {
	Api     ApiConfig
	Ssh     SshConfig
	Debug   bool
	Mongodb MongodbConfig
}

type MongodbConfig struct {
	Host string
	Db   string
}

type SshConfig struct {
	Listen         string
	PrivateKeyFile string
}

type ApiConfig struct {
	Listen string
}

var Cfg AppConfig

func Load(filename string) error {
	_, err := toml.DecodeFile(filename, &Cfg)
	return err
}
