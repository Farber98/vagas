package models

type Cards struct {
	IdCard uint64 `json:"id_card,omitempty"`
	*CardTypes
	Number     string `json:"card_number,omitempty"`
	Holder     string `json:"card_holder,omitempty"`
	Cvv        string `json:"cvv,omitempty"`
	ExpireDate string `json:"expire_date,omitempty"`
}

type CardTypes struct {
	IdCardType uint8  `json:"id_card_type,omitempty"`
	CardType   string `json:"card_type,omitempty"`
}
