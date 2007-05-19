package services

type ListTokenArgs struct {
	OwnerAddress     string
	TokenId          string
	Nonce            uint64
	TokenName        string
	FirstLink        string
	LastLink         string
	Hash             string
	Attributes       string
	Price            string
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
}

type BuyTokenArgs struct {
	OwnerAddress string
	BuyerAddress string
	TokenId      string
	Nonce        uint64
	Price        string
	Timestamp    uint64
	TxHash       string
}

type WithdrawTokenArgs struct {
	OwnerAddress string
	TokenId      string
	Nonce        uint64
	Price        string
	Timestamp    uint64
	TxHash       string
}
