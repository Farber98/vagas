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

func (srv *CardsService) ListCardTypes() ([]*models.CardTypes, error) {
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
