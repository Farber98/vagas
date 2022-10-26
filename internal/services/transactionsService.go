package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"pagarme/internal/constants"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
	"strconv"
	"time"
)

type TransactionsService struct {
	Db *infraestructure.DbHandler
}

func (srv *TransactionsService) Validate(tx *models.Transactions) error {

	if tx.IdClient <= 0 {
		return errors.New(constants.ERR_ID_CLIENT)
	}

	if tx.Value == "" {
		return errors.New(constants.ERR_VALUE)
	}

	intValue, err := strconv.Atoi(tx.Value)
	if err != nil {
		return errors.New(constants.ERR_VALUE)
	}

	if intValue <= 0 {
		return errors.New(constants.ERR_VALUE)
	}

	if tx.PaymentMethod != constants.CARD_TYPE_DEBIT && tx.PaymentMethod != constants.CARD_TYPE_CREDIT {
		return errors.New(constants.ERR_PAYMENT_METHOD)
	}

	if tx.Cards.Number == "" || len(tx.Cards.Number) != 16 {
		return errors.New(constants.ERR_CARD_NUMBER)
	}

	intCardNumber, err := strconv.Atoi(tx.Number)
	if err != nil {
		return errors.New(constants.ERR_CARD_NUMBER)
	}

	if intCardNumber < 1000000000000000 || intCardNumber > 9999999999999999 {
		return errors.New(constants.ERR_CARD_NUMBER)
	}

	if tx.Cards.Holder == "" {
		return errors.New(constants.ERR_CARD_HOLDER)
	}

	if tx.ExpireDate == "" {
		return errors.New(constants.ERR_CARD_DATE)
	}

	date, err := time.Parse(constants.DATE_LAYOUT, tx.ExpireDate)
	if err != nil {
		return errors.New(constants.ERR_CARD_DATE)
	}

	if date.Before(time.Now()) {
		return errors.New(constants.ERR_CARD_DATE)
	}

	if tx.Cvv == "" || len(tx.Cvv) != 3 {
		return errors.New(constants.ERR_CARD_CVV)
	}

	intCvv, err := strconv.Atoi(tx.Cvv)
	if err != nil {
		return errors.New(constants.ERR_CARD_CVV)
	}

	if intCvv < 100 || intCvv > 999 {
		return errors.New(constants.ERR_CARD_CVV)
	}

	return nil
}

func (srv *TransactionsService) Create(tx *models.Transactions) (*models.Transactions, error) {
	jsonPayload, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_transactions_create(?)", string(jsonPayload))

	var out infraestructure.SpOut

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, tx); err != nil {
		return nil, err
	}

	return tx, nil
}
