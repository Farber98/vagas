package controllers

import (
	"net/http"
	"pagarme/internal/models"
	"pagarme/internal/services"

	"github.com/labstack/echo/v4"
)

type ClientsController struct {
	*services.ClientsService
	*services.CardsService
}

func (controller *ClientsController) LoadRoutes(gr *echo.Group) {
	gr.GET("/hello", controller.Hello)
}

func (controller *ClientsController) Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, models.NewMsgResponse(OK_HELLO))
}
