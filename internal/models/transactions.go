package models

import "time"

type Transactions struct {
	IdTx        uint64    `json:"id_tx,omitempty"`
	IdCard      uint64    `json:"id_card,omitempty"`
	IdClient    uint32    `json:"id_client,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Value       string    `json:"value,omitempty"`
	Fee         string    `json:"fee,omitempty"`
	Status      string    `json:"status,omitempty"`
}
