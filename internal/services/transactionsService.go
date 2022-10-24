package services

import (
	"database/sql"
	"encoding/json"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
)

type TransactionsService struct {
	Db *infraestructure.DbHandler
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
