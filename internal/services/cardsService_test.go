package services

import (
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	CARD_TYPE_DEBIT  = "debit"
	CARD_TYPE_CREDIT = "credit"
)

func TestListCardTypes(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &CardsService{Db: db}

	t.Run("SUCCESS_DEBIT_AND_CREDIT_RETURNED", func(t *testing.T) {
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
}

func TestCreate(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &CardsService{Db: db}

	t.Run("SUCCESS", func(t *testing.T) {

		input := &models.Cards{
			CardTypes: models.CardTypes{
				IdCardType: 1,
				CardType:   CARD_TYPE_DEBIT,
			},
			Number:     "1111222233334444",
			Holder:     "Debit Holder",
			Cvv:        "123",
			ExpireDate: "2022-11-10",
		}

		expectedOutput := &models.Cards{
			IdCard: 1,
			CardTypes: models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     "1111222233334444",
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.Create(input)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("ERR_ALREADY_EXISTS", func(t *testing.T) {

		input := &models.Cards{
			CardTypes: models.CardTypes{
				IdCardType: 1,
			},
			Number:     "1111222233334444",
			Holder:     "Debit Holder",
			Cvv:        "123",
			ExpireDate: "2022-11-10",
		}

		expectedError := &mysql.MySQLError{
			Number:  0x426,
			Message: "Duplicate entry '1111222233334444' for key 'cards.ui_number'",
		}

		_, err := service.Create(input)
		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), expectedError.Error())
		}
	})
}

func TestFetchCard(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &CardsService{Db: db}

	t.Run("SUCCESS", func(t *testing.T) {
		idCard := uint64(1)

		expectedOutput := &models.Cards{
			IdCard: 1,
			CardTypes: models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     "1111222233334444",
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.Fetch(idCard)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("NOT_EXISTS", func(t *testing.T) {
		idCard := uint64(15)

		var expectedOutput *models.Cards

		card, err := service.Fetch(idCard)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})
}
func TestFetchCardByNumber(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	service := &CardsService{Db: db}

	t.Run("SUCCESS", func(t *testing.T) {
		number := "1111222233334444"

		expectedOutput := &models.Cards{
			IdCard: 1,
			CardTypes: models.CardTypes{
				CardType: CARD_TYPE_DEBIT,
			},
			Number:     "1111222233334444",
			Holder:     "Debit Holder",
			ExpireDate: "2022-11-10",
		}

		card, err := service.FetchByNumber(number)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})

	t.Run("NOT_EXISTS", func(t *testing.T) {
		number := "9999999999999999"

		var expectedOutput *models.Cards

		card, err := service.FetchByNumber(number)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, card)
		}
	})
}
