package config

import "errors"

// DBConfig Database config model
type DBConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Schema   string `toml:"schema"`
}

func (conf *DBConfig) SelfCheck() error {
	if conf.Host == "" {
		return errors.New("invalid host")
	}
	if conf.Port == "" {
		return errors.New("invalid port")
	}
	if conf.Username == "" {
		return errors.New("invalid username")
	}
	if conf.Password == "" {
		return errors.New("invalid password")
	}
	if conf.Schema == "" {
		return errors.New("invalid schema")
	}
	return nil
}
