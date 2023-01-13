package models

type Clients struct {
	IdClient uint32 `json:"id_client,omitempty" query:"id_client"`
	Wallets
}

type Wallets struct {
	IdWallet       uint32 `json:"id_wallet,omitempty"`
	AvailableFunds string `json:"available_funds,omitempty"`
	WaitingFunds   string `json:"waiting_funds,omitempty"`
}
