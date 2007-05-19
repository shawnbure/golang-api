package formatter

import (
	"encoding/hex"
	"fmt"
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
	mintTokensEndpointName      = "mintTokens"
	makeOfferEndpointName       = "makeOffer"
	acceptOfferEndpointName     = "acceptOffer"
	cancelOfferEndpointName     = "cancelOffer"
	startAuctionEndpointName    = "startAuction"
	placeBidEndpointName        = "placeBid"
	endAuctionEndpointName      = "endAuction"
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
		RcvAddr:   senderAddr,
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

func (f *TxFormatter) NewBuyNftTxTemplate(senderAddr string, tokenId string, nonce uint64, price string) Transaction {
	txData := buyNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     price,
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

func (f *TxFormatter) MakeOfferTxTemplate(senderAddr string, tokenId string, nonce uint64, amount float64) Transaction {
	txData := makeOfferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(amount).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.MakeOfferGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) AcceptOfferTxTemplate(senderAddr string, tokenId string, nonce uint64, offeror string, amount float64) (*Transaction, error) {
	offerorAddress, err := data.NewAddressFromBech32String(offeror)
	if err != nil {
		return nil, err
	}

	txData := makeOfferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(offerorAddress.AddressBytes()) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(amount).Bytes())

	return &Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.AcceptOfferGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}

func (f *TxFormatter) CancelOfferTxTemplate(senderAddr string, tokenId string, nonce uint64, amount float64) Transaction {
	txData := makeOfferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(amount).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.CancelOfferGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) NewMintNftsTxTemplate(
	walletAddress string,
	contractAddress string,
	mintPricePerToken float64,
	numberOfTokens uint64,
) Transaction {
	gasLimit := f.config.MintTokenGasLimit * numberOfTokens
	totalPrice := fmt.Sprintf("%f", mintPricePerToken*float64(numberOfTokens))
	txData := mintTokensEndpointName +
		"@" + hex.EncodeToString(big.NewInt(int64(numberOfTokens)).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     totalPrice,
		RcvAddr:   contractAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  gasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}
