package services

import (
	"database/sql"
	"encoding/json"
	infraestructure "pagarme/internal/infraestructures"
	"pagarme/internal/models"
)

type CardsService struct {
	Db *infraestructure.DbHandler
}

func (srv *CardsService) ListTypes() ([]*models.CardTypes, error) {
	row := srv.Db.QueryRow("CALL pg_card_list_types()")

	var out infraestructure.SpOut
	cardTypes := make([]*models.CardTypes, 0)
	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, &cardTypes); err != nil {
		return nil, err
	}

	return cardTypes, nil
}

func (srv *CardsService) Create(card *models.Cards) (*models.Cards, error) {
	jsonPayload, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_card_create(?)", string(jsonPayload))

	var out infraestructure.SpOut
	dbCard := &models.Cards{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, dbCard); err != nil {
		return nil, err
	}

	return dbCard, nil
}

func (srv *CardsService) Fetch(idCard uint64) (*models.Cards, error) {
	search := models.Search{}
	search["id_card"] = idCard

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_card_fetch(?)", string(jsonPayload))

	var out infraestructure.SpOut
	dbCard := &models.Cards{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, dbCard); err != nil {
		return nil, err
	}

	return dbCard, nil
}

func (srv *CardsService) FetchByNumber(number string) (*models.Cards, error) {
	search := models.Search{}
	search["card_number"] = number

	jsonPayload, err := json.Marshal(search)
	if err != nil {
		return nil, err
	}

	row := srv.Db.QueryRow("CALL pg_card_fetch_by_number(?)", string(jsonPayload))

	var out infraestructure.SpOut
	dbCard := &models.Cards{}

	if err := row.Scan(&out); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	if err := json.Unmarshal(out, dbCard); err != nil {
		return nil, err
	}

	return dbCard, nil
}
