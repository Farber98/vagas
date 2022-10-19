package config

import (
	"errors"
	"strconv"
)

// ContextConfig Context config model
type ContextConfig struct {
	Host  string `toml:"host"`
	Port  string `toml:"port"`
	Debug bool   `toml:"debug"`
}

func (conf *ContextConfig) SelfCheck() error {

	if conf.Host == "" {
		return errors.New("invalid hostname")
	}
	if _, err := strconv.ParseInt(conf.Port, 10, 64); conf.Port == "" || err != nil {
		return errors.New("invalid port")
	}

	return nil
}
