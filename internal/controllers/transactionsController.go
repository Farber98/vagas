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
	reqTx := &models.Transactions{}

	if err := c.Bind(reqTx); err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_BINDING))
	}

	err := controller.TransactionsService.Validate(reqTx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewMsgResponse(err.Error()))
	}

	client, err := controller.ClientsService.Fetch(reqTx.IdClient)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
	}

	if client == nil {
		client = &models.Clients{IdClient: reqTx.IdClient}
		client, err = controller.ClientsService.Create(client)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	card, err := controller.CardsService.FetchByNumber(reqTx.Cards.Number)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
	}

	if card == nil {
		card = &models.Cards{
			CardTypes: &models.CardTypes{
				IdCardType: uint8(dictionaries.Get().CardTypes[reqTx.PaymentMethod]),
			},
			Number:     reqTx.Cards.Number,
			Cvv:        reqTx.Cards.Cvv,
			Holder:     reqTx.Cards.Holder,
			ExpireDate: reqTx.Cards.ExpireDate,
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

	tx := &models.Transactions{
		IdClient: client.IdClient,
		Cards: &models.Cards{
			IdCard: card.IdCard,
		},
		Value:       reqTx.Value,
		Description: reqTx.Description,
	}

	tx, err = controller.TransactionsService.Create(tx)
	if err != nil {
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewMsgResponse(constants.ERR_DEFAULT))
		}
	}

	return c.JSON(http.StatusOK, models.NewDataResponse(constants.OK_TX, tx))
}
