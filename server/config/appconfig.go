package appconfig

import (
	"github.com/BurntSushi/toml"
)

/**
 * Config
 */
type AppConfig struct {
	Listen  string
	Debug   bool
	Mongodb MongodbConfig
}

type MongodbConfig struct {
	Host string
	Db   string
}

var Cfg AppConfig

func Load(filename string) error {
	_, err := toml.DecodeFile(filename, &Cfg)
	return err
}
