package services

import (
	"fmt"
	"math/rand"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	CARD_TYPE_DEBIT  = "debit"
	CARD_TYPE_CREDIT = "credit"
)

func TestCardsService(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &CardsService{Db: db}

	rand.Seed(time.Now().UnixNano())
	randCardNumber := fmt.Sprintf("%16d", rand.Int63n(1e16))
	wrongCardNumber := "9999"
	wrongIdCard := uint64(9999)

	var createdIdCard uint64

	t.Run("OK_LIST_CARD_TYPES", func(t *testing.T) {
		expectedOutput := []*models.CardTypes{
			{
				IdCardType: 1,
				CardType:   CARD_TYPE_DEBIT,
			},
			{
				IdCardType: 2,
				CardType:   CARD_TYPE_CREDIT,
			},
		}

		cardTypes, err := service.ListTypes()
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, cardTypes)
		}
	})

	t.Run("OK_CREATE", func(t *testing.T) {

		input := &models.Cards{
			CardTypes: &models.CardTypes{
				IdCardType: 1,
				CardType:   CARD_TYPE_DEBIT,
			},
			Number:     randCardNumber,
			Holder:     "Debit Holder",
			Cvv:        "123",
			ExpireDate: "2022-11-10",
		}

		expectedOutput := &models.Cards{
			CardTypes: &models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     randCardNumber,
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.Create(input)
		if assert.Nil(t, err) {
			expectedOutput.IdCard = card.IdCard
			createdIdCard = card.IdCard
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("ERR_CREATE_NUMBER_ALREADY_EXISTS", func(t *testing.T) {

		input := &models.Cards{
			CardTypes: &models.CardTypes{
				IdCardType: 1,
			},
			Number:     randCardNumber,
			Holder:     "Debit Holder",
			Cvv:        "123",
			ExpireDate: "2022-11-10",
		}

		expectedError := &mysql.MySQLError{
			Number:  0x426,
			Message: "Duplicate entry '" + randCardNumber + "' for key 'cards.ui_number'",
		}

		_, err := service.Create(input)
		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), expectedError.Error())
		}
	})

	t.Run("OK_FETCH", func(t *testing.T) {
		expectedOutput := &models.Cards{
			IdCard: createdIdCard,
			CardTypes: &models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     randCardNumber,
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.Fetch(createdIdCard)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("ERR_NOT_EXISTS_FETCH", func(t *testing.T) {
		var expectedOutput *models.Cards

		card, err := service.Fetch(wrongIdCard)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("OK_FETCH_BY_NUMBER", func(t *testing.T) {
		expectedOutput := &models.Cards{
			IdCard: createdIdCard,
			CardTypes: &models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     randCardNumber,
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.FetchByNumber(randCardNumber)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("ERR_NOT_EXISTS_FETCH_BY_NUMBER", func(t *testing.T) {
		var expectedOutput *models.Cards

		card, err := service.FetchByNumber(wrongCardNumber)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})
}
