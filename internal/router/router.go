package router

import (
	"log"
	"pagarme/internal/controllers"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/interfaces"
	"pagarme/internal/services"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var onceEcho sync.Once
var echoInstance *echo.Echo

//InitRoutes Initializes API routes.
func Init(db *infraestructure.DbHandler) *echo.Echo {
	onceEcho.Do(func() {
		e := echo.New()
		e.Use(middleware.CORS())
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogMethod:  true,
			LogURI:     true,
			LogStatus:  true,
			LogLatency: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				log.Printf("%v %v | Status: %v | Latency: %v\n", v.Method, v.URI, v.Status, v.Latency)
				return nil
			},
		}))

		/* Services */
		cardsService := &services.CardsService{Db: db}
		clientsService := &services.ClientsService{Db: db}
		txsService := &services.TransactionsService{Db: db}

		/* Controllers */
		arrayControllers := make([]interfaces.IController, 0)
		arrayControllers = append(arrayControllers, &controllers.HelloController{})
		arrayControllers = append(arrayControllers, &controllers.ClientsController{ClientsService: clientsService, CardsService: cardsService})
		arrayControllers = append(arrayControllers, &controllers.TransactionsController{ClientsService: clientsService, TransactionsService: txsService, CardsService: cardsService})

		group := e.Group("")
		for _, c := range arrayControllers {
			c.LoadRoutes(group)
		}

		echoInstance = e
	})
	return echoInstance
}
