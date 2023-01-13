package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pagarme/internal/constants"
	dictionaries "pagarme/internal/dictionary"
	"pagarme/internal/generators"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"pagarme/internal/services"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTx(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	dictionaries.Init(db)

	transactionsService := &services.TransactionsService{Db: db}
	clientsService := &services.ClientsService{Db: db}
	cardsService := &services.CardsService{Db: db}

	controller := &TransactionsController{
		ClientsService:      clientsService,
		TransactionsService: transactionsService,
		CardsService:        cardsService,
	}

	randIdClient := uint32(generators.RandomInt32(1, 999999))
	randDebitCard := fmt.Sprintf("%16d", generators.RandomInt64(1111111111111111, 9999999999999999))
	randCreditCard := fmt.Sprintf("%16d", generators.RandomInt64(1111111111111111, 9999999999999999))

	debitTransaction := &models.Transactions{
		IdClient:      randIdClient,
		PaymentMethod: constants.CARD_TYPE_DEBIT,
		Cards: &models.Cards{
			Number:     randDebitCard,
			Holder:     "The Peps",
			Cvv:        fmt.Sprintf("%3d", generators.RandomInt64(100, 999)),
			ExpireDate: "2030-11-10",
		},
		Value:       "100",
		Description: "OK_DEBIT_CLIENT_CARD_NOT_EXISTED",
	}
	creditTransaction := &models.Transactions{
		IdClient:      randIdClient,
		PaymentMethod: constants.CARD_TYPE_CREDIT,
		Cards: &models.Cards{
			Number:     randCreditCard,
			Holder:     "The Peps",
			Cvv:        fmt.Sprintf("%3d", generators.RandomInt64(100, 999)),
			ExpireDate: "2030-11-10",
		},
		Value:       "100",
		Description: "OK_CREDIT_ALL_EXISTED",
	}

	debitTransactionBytes, _ := json.Marshal(debitTransaction)
	creditTransactionBytes, _ := json.Marshal(creditTransaction)

	testCases := []struct {
		name             string
		body             string
		expectStatusCode int
		expectError      string
		outputTx         *models.Transactions
	}{
		{
			name:             "ERR_BINDING",
			body:             `NOT A JSON`,
			expectStatusCode: http.StatusInternalServerError,
			expectError:      constants.ERR_BINDING,
		},
		{
			name:             "ERR_VALIDATE_EMPTY_OBJECT",
			body:             `{}`,
			expectStatusCode: http.StatusBadRequest,
			expectError:      constants.ERR_ID_CLIENT,
		},
		{
			name:             "OK_DEBIT_CLIENT_CARD_NOT_EXISTED",
			body:             string(debitTransactionBytes),
			expectStatusCode: http.StatusOK,
			outputTx: &models.Transactions{
				IdClient: randIdClient,
				Cards: &models.Cards{
					CardTypes: &models.CardTypes{
						CardType: constants.CARD_TYPE_DEBIT,
					},
					Number: strings.Repeat("*", 12) + randDebitCard[12:],
					Holder: "The Peps",
				},
				Value:       "100.00",
				Fee:         "3.00",
				Description: "OK_DEBIT_CLIENT_CARD_NOT_EXISTED",
				Status:      constants.TX_STATUS_PAID,
			},
		},
		{
			name:             "OK_CREDIT_ALL_EXISTED",
			body:             string(creditTransactionBytes),
			expectStatusCode: http.StatusOK,
			outputTx: &models.Transactions{
				IdClient: randIdClient,
				Cards: &models.Cards{
					CardTypes: &models.CardTypes{
						CardType: constants.CARD_TYPE_CREDIT,
					},
					Number: strings.Repeat("*", 12) + randCreditCard[12:],
					Holder: "The Peps",
				},
				Value:       "100.00",
				Fee:         "5.00",
				Description: "OK_CREDIT_ALL_EXISTED",
				Status:      constants.TX_STATUS_WAITING_FUNDS,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/tx/create", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.Create(c)) {

				response := &models.MsgResponse{}
				err := json.Unmarshal(rec.Body.Bytes(), response)

				if assert.NoError(t, err) {
					assert.Equal(t, tc.expectStatusCode, rec.Code)
					if rec.Code != http.StatusOK {
						assert.Equal(t, tc.expectError, response.Message)
					} else {
						bytesRespTx, _ := json.Marshal(response.Data)
						respTx := &models.Transactions{}
						json.Unmarshal(bytesRespTx, respTx)

						tc.outputTx.IdCard = respTx.IdCard
						tc.outputTx.IdTx = respTx.IdTx
						tc.outputTx.Date = respTx.Date

						assert.Equal(t, respTx, tc.outputTx)
					}

				}
			}
		})
	}
}
