package formatter

import (
	"encoding/hex"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/services"
	"math/big"
)

var (
	listNftEndpointName         = "putNftForSale"
	buyNftEndpointName          = "buyNft"
	withdrawNftEndpointName     = "withdrawNft"
	ESDTNFTTransferEndpointName = "ESDTNFTTransfer"
)

type TxFormatter struct {
	config config.BlockchainConfig
}

func NewTxFormatter(cfg config.BlockchainConfig) TxFormatter {
	return TxFormatter{config: cfg}
}

func (f *TxFormatter) NewListNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price float64) (*data.Transaction, error) {
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

	tx := data.Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.ListNftGasLimit,
		Data:      []byte(txData),
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}

	return &tx, nil
}

func (f *TxFormatter) NewBuyNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price float64) data.Transaction {
	priceDecStr := services.GetPriceDenominated(price).Text(10)

	txData := buyNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return data.Transaction{
		Nonce:     0,
		Value:     priceDecStr,
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.BuyNftGasLimit,
		Data:      []byte(txData),
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) NewWithdrawNftTxTemplate(senderAddr string, tokenId string, nonce uint64) data.Transaction {
	txData := withdrawNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return data.Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.WithdrawNftGasLimit,
		Data:      []byte(txData),
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}
