package services

type ListTokenArgs struct {
	OwnerAddress     string
	BuyerAddress     string
	TokenId          string
	Nonce            uint64
	NonceStr         string
	TokenName        string
	FirstLink        string
	SecondLink       string
	Hash             string
	Attributes       string
	Price            string
	PriceNominal     string
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
	TxConfirmed      bool
	OnSale           bool
	AuctionStartTime uint64
	AuctionDeadline  uint64
}

type BuyTokenArgs struct {
	OwnerAddress string
	BuyerAddress string
	TokenId      string
	Nonce        uint64
	NonceStr     string
	Price        string
	PriceNominal string
	Timestamp    uint64
	TxHash       string
	TxConfirmed  bool
	OnSale       bool
}

type WithdrawTokenArgs struct {
	OwnerAddress string
	TokenId      string
	Nonce        uint64
	NonceStr     string
	Price        string
	PriceNominal string
	Timestamp    uint64
	TxHash       string
	TxConfirmed  bool
	OnSale       bool
}
