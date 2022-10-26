package services

import (
	"fmt"
	"math/rand"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TX_STATUS_PAID    = "P"
	TX_STATUS_WAITING = "W"
)

func TestCreateTx(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &TransactionsService{Db: db}
	clientService := &ClientsService{Db: db}
	cardService := &CardsService{Db: db}

	rand.Seed(time.Now().UnixNano())
	randDebitCardNumber := fmt.Sprintf("%16d", rand.Int63n(1e16))
	randCreditCardNumber := fmt.Sprintf("%16d", rand.Int63n(1e16))
	randIdClient := rand.Int31n(1000)

	//Create DEBIT card.
	debitCard := &models.Cards{
		CardTypes: &models.CardTypes{
			IdCardType: 1,
			CardType:   CARD_TYPE_DEBIT,
		},
		Number:     randDebitCardNumber,
		Holder:     "Debit Holder",
		Cvv:        "123",
		ExpireDate: "2022-11-10",
	}

	createdDebitCard, _ := cardService.Create(debitCard)

	//Create CREDIT card.
	creditCard := &models.Cards{
		CardTypes: &models.CardTypes{
			IdCardType: 2,
			CardType:   CARD_TYPE_CREDIT,
		},
		Number:     randCreditCardNumber,
		Holder:     "Credit Holder",
		Cvv:        "123",
		ExpireDate: "2022-11-10",
	}

	createdCreditCard, _ := cardService.Create(creditCard)

	// Create client
	clientService.Create(&models.Clients{IdClient: uint32(randIdClient)})

	t.Run("SUCCESS_DEBIT", func(t *testing.T) {

		input := &models.Transactions{
			IdClient: uint32(randIdClient),
			Cards: &models.Cards{
				IdCard: createdDebitCard.IdCard,
			},
			Value:       "100",
			Description: "Smartband XYZ 1.0",
		}

		expectedOutput := &models.Transactions{
			IdClient: uint32(randIdClient),
			Cards: &models.Cards{
				IdCard: createdDebitCard.IdCard,
				CardTypes: &models.CardTypes{
					CardType: CARD_TYPE_DEBIT,
				},
				Holder: "Debit Holder",
				Number: strings.Repeat("*", 12) + randDebitCardNumber[12:],
			},
			Value:       "100.00",
			Fee:         "3.00",
			Description: "Smartband XYZ 1.0",
			Status:      "P",
		}

		tx, err := service.Create(input)

		if assert.Nil(t, err) {
			expectedOutput.IdTx = tx.IdTx
			expectedOutput.Date = tx.Date
			assert.Equal(t, expectedOutput, tx)
		}
	})

	t.Run("SUCCESS_CREDIT", func(t *testing.T) {

		input := &models.Transactions{
			IdClient: uint32(randIdClient),
			Cards: &models.Cards{
				IdCard: createdCreditCard.IdCard,
			},
			Value:       "100",
			Description: "Smartband XYZ 2.0",
		}

		expectedOutput := &models.Transactions{
			IdClient: uint32(randIdClient),
			Cards: &models.Cards{
				IdCard: createdCreditCard.IdCard,
				CardTypes: &models.CardTypes{
					CardType: CARD_TYPE_CREDIT,
				},
				Holder: "Credit Holder",
				Number: strings.Repeat("*", 12) + randCreditCardNumber[12:],
			},
			Value:       "100.00",
			Fee:         "5.00",
			Description: "Smartband XYZ 2.0",
			Status:      "W",
		}

		tx, err := service.Create(input)
		if assert.Nil(t, err) {
			expectedOutput.IdTx = tx.IdTx
			expectedOutput.Date = tx.Date
			assert.Equal(t, expectedOutput, tx)
		}
	})

}
