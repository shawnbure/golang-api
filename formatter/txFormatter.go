package formatter

import (
	"encoding/hex"
	"math/big"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/services"
)

var (
	listNftEndpointName         = "putNftForSale"
	buyNftEndpointName          = "buyNft"
	withdrawNftEndpointName     = "withdrawNft"
	ESDTNFTTransferEndpointName = "ESDTNFTTransfer"
)

type Transaction struct {
	Nonce     uint64 `json:"nonce"`
	Value     string `json:"value"`
	RcvAddr   string `json:"receiver"`
	SndAddr   string `json:"sender"`
	GasPrice  uint64 `json:"gasPrice,omitempty"`
	GasLimit  uint64 `json:"gasLimit,omitempty"`
	Data      string `json:"data,omitempty"`
	Signature string `json:"signature,omitempty"`
	ChainID   string `json:"chainID"`
	Version   uint32 `json:"version"`
	Options   uint32 `json:"options,omitempty"`
}

type TxFormatter struct {
	config config.BlockchainConfig
}

func NewTxFormatter(cfg config.BlockchainConfig) TxFormatter {
	return TxFormatter{config: cfg}
}

func (f *TxFormatter) NewListNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price float64) (*Transaction, error) {
	marketPlaceAddress, err := data.NewAddressFromBech32String(f.config.MarketplaceAddress)
	if err != nil {
		return nil, err
	}

	txData := ESDTNFTTransferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(1)).Bytes()) +
		"@" + hex.EncodeToString(marketPlaceAddress.AddressBytes()) +
		"@" + hex.EncodeToString([]byte(listNftEndpointName)) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(price).Bytes())

	tx := Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.ListNftGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}

	return &tx, nil
}

func (f *TxFormatter) NewBuyNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price float64) Transaction {
	priceDecStr := services.GetPriceDenominated(price).Text(10)

	txData := buyNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     priceDecStr,
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.BuyNftGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) NewWithdrawNftTxTemplate(senderAddr string, tokenId string, nonce uint64) Transaction {
	txData := withdrawNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.WithdrawNftGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}
