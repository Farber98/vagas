package main

import (
	"log"
	"pagarme/internal/config"
	"pagarme/internal/router"
)

func main() {
	e := router.Init()
	log.Fatal(e.Start(config.Get().Context.Host + ":" + config.Get().Context.Port))
}
