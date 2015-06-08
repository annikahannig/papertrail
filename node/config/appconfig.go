package appconfig

import (
	"github.com/BurntSushi/toml"
)

/**
 * Config
 */
type AppConfig struct {
	Server   string
	NodeName string
}

var Cfg AppConfig

func Load(filename string) error {
	_, err := toml.DecodeFile(filename, &Cfg)
	return err
}
