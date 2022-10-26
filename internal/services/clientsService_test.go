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

func TestClients(t *testing.T) {

	db := infraestructure.ConstructTestDB()
	clientService := &ClientsService{Db: db}
	cardService := &CardsService{Db: db}
	txService := &TransactionsService{Db: db}

	rand.Seed(time.Now().UnixNano())
	randCardNumber := fmt.Sprintf("%16d", rand.Int63n(1e16))
	randIdClient := rand.Int31n(1000)
	otherRandIdClient := rand.Int31n(1000)
	//randCardNumber := fmt.Sprintf("%16d", rand.Int63n(1e16))

	wrongIdClient := uint32(99999)

	var createdIdCard uint64
	var createdIdClient uint32

	//Create card.

	card := &models.Cards{
		CardTypes: &models.CardTypes{
			IdCardType: 1,
			CardType:   CARD_TYPE_DEBIT,
		},
		Number:     randCardNumber,
		Holder:     "Debit Holder",
		Cvv:        "123",
		ExpireDate: "2022-11-10",
	}
	createdCard, _ := cardService.Create(card)

	createdIdCard = createdCard.IdCard

	t.Run("SUCCESS", func(t *testing.T) {

		input := &models.Clients{
			IdClient: uint32(randIdClient),
		}

		expectedOutput := &models.Clients{
			IdClient: uint32(randIdClient),
			Wallets: models.Wallets{
				IdWallet:       uint32(randIdClient),
				AvailableFunds: "0.00",
				WaitingFunds:   "0.00",
			},
		}

		client, err := clientService.Create(input)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, client)
			createdIdClient = client.IdClient
		}
	})

	//Creates a new client and adds a transaction.
	clientService.Create(&models.Clients{IdClient: uint32(otherRandIdClient)})
	tx := &models.Transactions{
		IdClient: uint32(otherRandIdClient),
		Cards: &models.Cards{
			IdCard: createdIdCard,
		},
		Value:       "100",
		Description: "Smartband XYZ 1.0",
	}

	txService.Create(tx)

	t.Run("ERR_ALREADY_EXISTS", func(t *testing.T) {

		input := &models.Clients{
			IdClient: createdIdClient,
		}

		expectedError := &mysql.MySQLError{
			Number:  0x426,
			Message: "Duplicate entry '" + fmt.Sprint(createdIdClient) + "' for key 'clients.PRIMARY'",
		}

		_, err := clientService.Create(input)
		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), expectedError.Error())
		}
	})

	t.Run("SUCCESS_FETCH", func(t *testing.T) {
		expectedOutput := &models.Clients{
			IdClient: createdIdClient,
			Wallets: models.Wallets{
				IdWallet:       createdIdClient,
				AvailableFunds: "0.00",
				WaitingFunds:   "0.00",
			},
		}

		client, err := clientService.Fetch(createdIdClient)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, client)
		}
	})

	t.Run("NOT_EXISTS_FETCH", func(t *testing.T) {

		var expectedOutput *models.Clients

		client, err := clientService.Fetch(wrongIdClient)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, client)
		}
	})

	t.Run("SUCCESS_REGISTER_CARD", func(t *testing.T) {
		expectedOutput := &models.ClientsCards{
			IdClient: createdIdClient,
			IdCard:   createdIdCard,
		}

		clientCard, err := clientService.RegisterCard(createdIdCard, createdIdClient)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedOutput, clientCard)
		}
	})

	t.Run("ERR_ALREADY_EXISTS_REGISTER_CARD", func(t *testing.T) {

		expectedError := &mysql.MySQLError{
			Number:  0x426,
			Message: "Duplicate entry '" + fmt.Sprint(createdIdCard) + "-" + fmt.Sprint(createdIdClient) + "' for key 'clients_cards.PRIMARY'",
		}

		_, err := clientService.RegisterCard(createdIdCard, createdIdClient)
		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), expectedError.Error())
		}
	})

	t.Run("OK_TX_RETURNED", func(t *testing.T) {

		txs, err := clientService.ListTransactions(createdIdClient)
		if assert.Nil(t, err) {
			assert.NotNil(t, txs)
		}
	})

	t.Run("OK_TX_EMPTY", func(t *testing.T) {

		txs, err := clientService.ListTransactions(wrongIdClient)
		if assert.Nil(t, err) {
			assert.Empty(t, txs)
		}
	})
}
