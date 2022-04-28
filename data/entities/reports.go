package entities

type TopVolumeByAddress struct {
	FromTime string  `json:"from_time"`
	ToTime   string  `json:"to_time"`
	Address  string  `json:"address"`
	Volume   float64 `json:"volume"`
}

type VerifiedListingTransaction struct {
	TxId              uint64  `json:"txId"`
	TxType            string  `json:"txType"`
	TxHash            string  `json:"txHash"`
	TxPriceNominal    float64 `json:"txPriceNominal"`
	TxTimestamp       int64   `json:"txTimestamp"`
	TokenId           string  `json:"tokenId"`
	TokenName         string  `json:"tokenName"`
	TokenImageLink    string  `json:"tokenImageLink"`
	Address           string  `json:"address"`
	CollectionTokenId string  `json:"collectionTokenId"`
	CollectionName    string  `json:"collectionName"`
}
