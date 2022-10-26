package models

type Transactions struct {
	IdTx uint64 `json:"id_tx,omitempty"`
	*Cards
	IdClient    uint32 `json:"id_client,omitempty"`
	Description string `json:"description,omitempty"`
	Date        string `json:"date,omitempty"`
	Value       string `json:"value,omitempty"`
	Fee         string `json:"fee,omitempty"`
	Status      string `json:"status,omitempty"`
}
