package main

import (
	"log"
	"pagarme/internal/config"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/router"
)

func main() {
	db := infraestructure.ConstructDB()
	e := router.Init(db)
	log.Fatal(e.Start(config.Get().Context.Host + ":" + config.Get().Context.Port))
}
