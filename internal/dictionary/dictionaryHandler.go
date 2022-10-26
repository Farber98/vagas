package dictionaries

import (
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/services"
	"sync"
)

type Handler struct {
	CardTypes map[string]int
}

var handlerInstance *Handler
var onceHandler sync.Once

func Init(db *infraestructure.DbHandler) *Handler {
	onceHandler.Do(func() {
		cardService := &services.CardsService{Db: db}
		cardTypes, err := cardService.ListTypes()
		if err != nil {
			panic(err)
		}

		h := &Handler{
			CardTypes: make(map[string]int),
		}
		for _, t := range cardTypes {
			h.CardTypes[t.CardType] = int(t.IdCardType)
		}
		handlerInstance = h
	})
	return handlerInstance
}

func Get() *Handler {
	return handlerInstance
}
