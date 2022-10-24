package main

import (
	"log"
	"pagarme/internal/config"
	dictionaries "pagarme/internal/dictionary"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/router"
)

func main() {
	db := infraestructure.ConstructDB()
	e := router.Init(db)
	dictionaries.Init(db)
	log.Fatal(e.Start(config.Get().Context.Host + ":" + config.Get().Context.Port))
}
