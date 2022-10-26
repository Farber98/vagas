package services

import (
	"database/sql"
	"encoding/json"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
)

type ClientsService struct {
	Db *infraestructure.DbHandler
}

func (srv *ClientsService) Create(client *models.Clients) (*models.Clients, error) {
	jsonPayload, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_client_create(?)", string(jsonPayload))

	var out infraestructure.SpOut

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (srv *ClientsService) RegisterCard(idCard uint64, idClient uint32) (*models.ClientsCards, error) {
	search := models.Search{}
	search["id_card"] = idCard
	search["id_client"] = idClient

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_client_register_card(?)", string(jsonPayload))

	var out infraestructure.SpOut
	clientCard := &models.ClientsCards{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, clientCard); err != nil {
		return nil, err
	}

	return clientCard, nil
}

func (srv *ClientsService) Fetch(idClient uint32) (*models.Clients, error) {
	search := models.Search{}
	search["id_client"] = idClient

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_client_fetch(?)", string(jsonPayload))

	var out infraestructure.SpOut
	client := &models.Clients{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (srv *ClientsService) FetchCard(idClient uint32, idCard uint64) (*models.ClientsCards, error) {
	search := models.Search{}
	search["id_client"] = idClient
	search["id_card"] = idClient

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_client_fetch_card(?)", string(jsonPayload))

	var out infraestructure.SpOut
	client := &models.ClientsCards{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (srv *ClientsService) ListTransactions(idClient uint32) ([]*models.Transactions, error) {
	search := models.Search{}
	search["id_client"] = idClient

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_client_list_transactions(?)", string(jsonPayload))

	var out infraestructure.SpOut
	txs := make([]*models.Transactions, 0)

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}
