package services

import (
	"fmt"
	"math/rand"
	"pagarme/internal/constants"
	"pagarme/internal/generators"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateTx(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	txsService := &TransactionsService{Db: db}

	testCases := []struct {
		name           string
		input          *models.Transactions
		expectedOutput string
	}{
		{
			name:           "ERR_ID_CLIENT",
			input:          &models.Transactions{},
			expectedOutput: constants.ERR_ID_CLIENT,
		},
		{
			name: "ERR_EMPTY_VALUE",
			input: &models.Transactions{
				IdClient: 1,
			},
			expectedOutput: constants.ERR_VALUE,
		},
		{
			name: "ERR_WRONG_VALUE",
			input: &models.Transactions{
				IdClient: 1,
				Value:    "wrong",
			},
			expectedOutput: constants.ERR_VALUE,
		},
		{
			name: "ERR_OUTRANGED_VALUE",
			input: &models.Transactions{
				IdClient: 1,
				Value:    "-5",
			},
			expectedOutput: constants.ERR_VALUE,
		},
		{
			name: "ERR_WRONG_PAYMENT_METHOD",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: "wrong",
			},
			expectedOutput: constants.ERR_PAYMENT_METHOD,
		},
		{
			name: "ERR_EMPTY_NUMBER",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards:         &models.Cards{},
			},
			expectedOutput: constants.ERR_CARD_NUMBER,
		},
		{
			name: "ERR_LENGTH_NUMBER",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number: "wrong",
				},
			},
			expectedOutput: constants.ERR_CARD_NUMBER,
		},
		{
			name: "ERR_WRONG_NUMBER",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number: "wrongwrongwrongw",
				},
			},
			expectedOutput: constants.ERR_CARD_NUMBER,
		},
		{
			name: "ERR_OUTRANGED_NUMBER",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number: "1111111111",
				},
			},
			expectedOutput: constants.ERR_CARD_NUMBER,
		},
		{
			name: "ERR_EMPTY_HOLDER",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number: "3333333333333333",
				},
			},
			expectedOutput: constants.ERR_CARD_HOLDER,
		},
		{
			name: "ERR_EMPTY_DATE",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number: "3333333333333333",
					Holder: "Some Holder",
				},
			},
			expectedOutput: constants.ERR_CARD_DATE,
		},
		{
			name: "ERR_WRONG_DATE",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					ExpireDate: "wrong",
				},
			},
			expectedOutput: constants.ERR_CARD_DATE,
		},
		{
			name: "ERR_EXPIRED_DATE",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					ExpireDate: constants.DATE_LAYOUT,
				},
			},
			expectedOutput: constants.ERR_CARD_DATE,
		},
		{
			name: "ERR_EMPTY_CVV",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					ExpireDate: "2040-01-02 15:04:05",
				},
			},
			expectedOutput: constants.ERR_CARD_CVV,
		},
		{
			name: "ERR_LENGTH_CVV",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					Cvv:        "wrong",
					ExpireDate: "2030-01-02 15:04:05",
				},
			},
			expectedOutput: constants.ERR_CARD_CVV,
		},
		{
			name: "ERR_WRONG_CVV",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					Cvv:        "wro",
					ExpireDate: "2030-01-02 15:04:05",
				},
			},
			expectedOutput: constants.ERR_CARD_CVV,
		},
		{
			name: "ERR_OUTRANGED_CVV",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					Cvv:        "10",
					ExpireDate: "2030-01-02 15:04:05",
				},
			},
			expectedOutput: constants.ERR_CARD_CVV,
		},
		{
			name: "OK",
			input: &models.Transactions{
				IdClient:      1,
				Value:         "100",
				PaymentMethod: constants.CARD_TYPE_CREDIT,
				Cards: &models.Cards{
					Number:     "3333333333333333",
					Holder:     "Some Holder",
					Cvv:        "333",
					ExpireDate: "2030-01-02 15:04:05",
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			serviceError := txsService.Validate(tc.input)
			if tc.expectedOutput != "" {
				if assert.Error(t, serviceError) {
					assert.Equal(t, tc.expectedOutput, serviceError.Error())
				}
			} else {
				assert.Nil(t, serviceError)
			}
		})
	}
}
func TestCreateTx(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &TransactionsService{Db: db}
	clientService := &ClientsService{Db: db}
	cardService := &CardsService{Db: db}

	rand.Seed(time.Now().UnixNano())
	randDebitCardNumber := fmt.Sprintf("%16d", generators.RandomInt64(1111111111111111, 9999999999999999))
	randCreditCardNumber := fmt.Sprintf("%16d", generators.RandomInt64(1111111111111111, 9999999999999999))
	randIdClient := generators.RandomInt32(1, 499999)

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
			Status:      constants.TX_STATUS_PAID,
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
			Status:      constants.TX_STATUS_WAITING_FUNDS,
		}

		tx, err := service.Create(input)
		if assert.Nil(t, err) {
			expectedOutput.IdTx = tx.IdTx
			expectedOutput.Date = tx.Date
			assert.Equal(t, expectedOutput, tx)
		}
	})

}
