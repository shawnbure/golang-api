package formatter

import (
	"encoding/hex"
	"errors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/config"
	"math/big"
	"strconv"
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

func NewTxFormatter(cfg config.BlockchainConfig, ) TxFormatter {
	return TxFormatter{config: cfg}
}

func (f *TxFormatter) NewListNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price string) (*data.Transaction, error) {
	marketPlaceAddress, err := data.NewAddressFromBech32String(f.config.MarketplaceAddress)
	if err != nil {
		return nil, err
	}

	priceBigInt, success := big.NewInt(0).SetString(price, 10)
	if !success {
		return nil, errors.New("cannot parse price")
	}

	txData := ESDTNFTTransferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(1)).Bytes()) +
		"@" + hex.EncodeToString(marketPlaceAddress.AddressBytes()) +
		"@" + hex.EncodeToString([]byte(listNftEndpointName)) +
		"@" + hex.EncodeToString(priceBigInt.Bytes())

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

func (f *TxFormatter) NewBuyNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price string) data.Transaction {
	priceDec := big.NewInt(0).SetBytes([]byte(price)).Int64()
	priceDecStr := strconv.FormatInt(priceDec, 10)

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
