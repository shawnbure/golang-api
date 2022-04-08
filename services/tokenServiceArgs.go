package services

type ListTokenArgs struct {
	OwnerAddress     string
	BuyerAddress     string
	TokenId          string
	Nonce            uint64
	TokenName        string
	FirstLink        string
	SecondLink       string
	Hash             string
	Attributes       string
	Price            string
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
	TxConfirmed      bool
	OnSale           bool
}

type BuyTokenArgs struct {
	OwnerAddress string
	BuyerAddress string
	TokenId      string
	Nonce        uint64
	StrNonce     string
	Price        string
	NominalPrice string
	Timestamp    uint64
	TxHash       string
	TxConfirmed  bool
	OnSale       bool
}

type WithdrawTokenArgs struct {
	OwnerAddress string
	TokenId      string
	Nonce        uint64
	Price        string
	Timestamp    uint64
	TxHash       string
	TxConfirmed  bool
	OnSale       bool
}
