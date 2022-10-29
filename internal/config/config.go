package config

import (
	"errors"
	"log"
	"pagarme/internal/interfaces"
	"reflect"
	"sync"

	"github.com/BurntSushi/toml"
)

var singleConfigInstance *Config
var onceConfig sync.Once

func Init() *Config {
	onceConfig.Do(func() {
		log.Println("Loading config...")
		c := Config{}
		if _, err := toml.DecodeFile("./config_local.toml", &c); err != nil {
			if _, err := toml.DecodeFile("./config_docker.toml", &c); err != nil {
				panic(err)
			}
		}
		if err := c.Check(); err != nil {
			panic(err)
		}
		singleConfigInstance = &c
		log.Println("Finished loading config")
	})
	return singleConfigInstance
}

//Get Returns loaded configuration structure
func Get() *Config {
	return Init()
}

// Check runs SelfCheck() for each Config struct attribute.
func (c *Config) Check() error {
	selfType := reflect.ValueOf(c).Elem()

	fields := selfType.NumField()
	for i := 0; i < fields; i++ {
		field := selfType.Field(i)
		structField := field.Interface()
		if iFaced, ok := structField.(interfaces.IConfig); ok {
			if err := iFaced.SelfCheck(); err != nil {
				return errors.New(field.Type().String() + ": " + err.Error())
			}
		} else {
			return errors.New("Configuration property " + field.Type().String() + " doesn't implement IConfig interface")
		}
	}
	return nil
}

// Config Config Model
type Config struct {
	Context *ContextConfig `toml:"context"`
	DB      *DBConfig      `toml:"database"`
	DBTest  *DBConfig      `toml:"database_test"`
}
