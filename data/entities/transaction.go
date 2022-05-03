package entities

type Transaction struct {
	ID           uint64  `gorm:"primaryKey" json:"id"`
	Hash         string  `json:"hash" gorm:"index:,unique"`
	Type         TxType  `json:"type" `
	PriceNominal float64 `json:"priceNominal"`
	Timestamp    uint64  `json:"timestamp"`
	SellerID     uint64  `json:"sellerId"`
	BuyerID      uint64  `json:"buyerId"`
	TokenID      uint64  `json:"tokenId"`
	CollectionID uint64  `json:"collectionId"`
}

type TxType string

const (
	ListToken     TxType = "List"
	BuyToken      TxType = "Buy"
	WithdrawToken TxType = "Withdraw"
	AuctionToken  TxType = "Auction"
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
}

type Activity struct {
	TxId              uint64  `json:"txId"`
	TxType            string  `json:"txType"`
	TxHash            string  `json:"txHash"`
	TxPriceNominal    float64 `json:"txPriceNominal"`
	TxTimestamp       int64   `json:"txTimestamp"`
	TokenId           string  `json:"tokenId"`
	TokenName         string  `json:"tokenName"`
	TokenImageLink    string  `json:"tokenImageLink"`
	FromAddress       string  `json:"fromAddress"`
	ToAddress         string  `json:"toAddress"`
	CollectionId      string  `json:"collectionId"`
	CollectionTokenId string  `json:"collectionTokenId"`
	CollectionName    string  `json:"collectionName"`
}

type AggregatedTradeVolume struct {
	BuyVolume      float64 `json:"buy_volume"`
	WithdrawVolume float64 `json:"withdraw_volume"`
	ListVolume     float64 `json:"list_volume"`
}
