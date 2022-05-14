package entities

type Transaction struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Hash         string     `json:"hash" gorm:"index:,unique"`
	Type         TxType     `json:"type" `
	PriceNominal float64    `json:"priceNominal"`
	Timestamp    uint64     `json:"timestamp"`
	Seller       Account    `json:"seller"`
	SellerID     uint64     `json:"sellerId"`
	Buyer        Account    `json:"buyer"`
	BuyerID      uint64     `json:"buyerId"`
	TokenID      uint64     `json:"tokenId"`
	Token        Token      `json:"token"  `
	CollectionID uint64     `json:"collectionId"`
	Collection   Collection `json:"collection"`
}

type TxType string

const (
	ListToken     TxType = "List"
	BuyToken      TxType = "Buy"
	WithdrawToken TxType = "Withdraw"
	AuctionToken  TxType = "Auction"
	TxStake       TxType = "Stake"
	None          TxType = "None"
)

type TransactionDetail struct {
	TxId           uint64  `json:"txId"`
	TxType         string  `json:"txType"`
	TxHash         string  `json:"txHash"`
	TxPriceNominal float64 `json:"txPriceNominal"`
	TxTimestamp    int64   `json:"txTimestamp"`
	TokenId        string  `json:"tokenId"`
	TokenName      string  `json:"tokenName"`
	TokenImageLink string  `json:"tokenImageLink"`
	FromAddress    string  `json:"fromAddress"`
	ToAddress      string  `json:"toAddress"`
	ToId           int64   `json:"to_id"`
}
type Activity struct {
	Transaction  Transaction `json:"transaction" gorm:"embedded"`
	Token        Token       `json:"token"`
	TokenID      uint64      `json:"tokenId"`
	Collection   Collection  `json:"collection"  `
	CollectionID uint64      `json:"collectionId"`
}

type AggregatedTradeVolume struct {
	BuyVolume      float64 `json:"buy_volume"`
	WithdrawVolume float64 `json:"withdraw_volume"`
	ListVolume     float64 `json:"list_volume"`
}
