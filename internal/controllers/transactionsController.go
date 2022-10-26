package controllers

import (
	"net/http"
	"pagarme/internal/constants"
	dictionaries "pagarme/internal/dictionary"
	"pagarme/internal/models"
	"pagarme/internal/services"

	"github.com/labstack/echo/v4"
)

type TransactionsController struct {
	*services.ClientsService
	*services.TransactionsService
	*services.CardsService
}

func (controller *TransactionsController) LoadRoutes(gr *echo.Group) {
	gr.POST("/tx/create", controller.Create)
}

func (controller *TransactionsController) Create(c echo.Context) error {
	tx := &models.Transactions{}

	if err := c.Bind(tx); err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_BINDING))
	}

	err := controller.TransactionsService.Validate(tx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewMsgResponse(err.Error()))
	}

	client, err := controller.ClientsService.Fetch(tx.IdClient)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
	}

	if client == nil {
		client = &models.Clients{IdClient: tx.IdClient}
		client, err = controller.ClientsService.Create(client)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	card, err := controller.CardsService.FetchByNumber(tx.Cards.Number)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
	}

	if card == nil {
		card = &models.Cards{
			CardTypes: &models.CardTypes{
				IdCardType: uint8(dictionaries.Get().CardTypes[tx.PaymentMethod]),
			},
			Number:     tx.Cards.Number,
			Cvv:        tx.Cards.Cvv,
			Holder:     tx.Cards.Holder,
			ExpireDate: tx.Cards.ExpireDate,
		}

		card, err = controller.CardsService.Create(card)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	clientCard, err := controller.ClientsService.FetchCard(client.IdClient, card.IdCard)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
	}

	if clientCard == nil {
		clientCard, err = controller.ClientsService.RegisterCard(card.IdCard, client.IdClient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	tx.Cards = card
	tx.IdClient = client.IdClient

	tx, err = controller.TransactionsService.Create(tx)
	if err != nil {
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	return c.JSON(http.StatusOK, models.NewDataResponse(constants.OK_TX, tx))
}
