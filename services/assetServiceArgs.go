package services

type ListAssetArgs struct {
	OwnerAddress     string
	TokenId          string
	Nonce            uint64
	TokenName        string
	FirstLink        string
	LastLink         string
	Hash             string
	Attributes       string
	Uri              string
	Price            string
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
}

type BuyAssetArgs struct {
	OwnerAddress string
	BuyerAddress string
	TokenId      string
	Nonce        uint64
	Uri          string
	Price        string
	Timestamp    uint64
	TxHash       string
}

type WithdrawAssetArgs struct {
	OwnerAddress string
	TokenId      string
	Nonce        uint64
	Uri          string
	Price        string
	Timestamp    uint64
	TxHash       string
}
