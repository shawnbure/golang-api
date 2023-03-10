package services

type MakeOfferArgs struct {
	OfferorAddress string
	TokenId        string
	Nonce          uint64
	Amount         string
	Expire         uint64
	Timestamp      uint64
	TxHash         string
}

type CancelOfferArgs struct {
	OfferorAddress string
	TokenId        string
	Nonce          uint64
	Amount         string
	Timestamp      uint64
	TxHash         string
}

type AcceptOfferArgs struct {
	OwnerAddress   string
	TokenId        string
	Nonce          uint64
	OfferorAddress string
	Amount         string
	Timestamp      uint64
	TxHash         string
}

type StartAuctionArgs struct {
	OwnerAddress     string
	TokenId          string
	Nonce            uint64
	TokenName        string
	FirstLink        string
	SecondLink       string
	Hash             string
	Attributes       string
	MinBid           string
	StartTime        uint64
	Deadline         uint64
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
}

type PlaceBidArgs struct {
	Offeror   string
	TokenId   string
	Nonce     uint64
	Amount    string
	Timestamp uint64
	TxHash    string
}

type EndAuctionArgs struct {
	Caller    string
	TokenId   string
	Nonce     uint64
	Winner    string
	Amount    string
	Timestamp uint64
	TxHash    string
}

type DepositUpdateArgs struct {
	Owner  string
	Amount string
}
