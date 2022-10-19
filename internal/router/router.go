package router

import (
	"log"
	"pagarme/internal/controllers"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/interfaces"
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

		arrayControllers := make([]interfaces.IController, 0)
		arrayControllers = append(arrayControllers, &controllers.HelloController{})

		group := e.Group("")
		for _, c := range arrayControllers {
			c.LoadRoutes(group)
		}

		echoInstance = e
	})
	return echoInstance
}
