package formatter

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"gorm.io/gorm"
)

var (
	listNftEndpointName                      = "putNftForSale"
	buyNftEndpointName                       = "buyNft"
	withdrawNftEndpointName                  = "withdrawNft"
	ESDTNFTTransferEndpointName              = "ESDTNFTTransfer"
	mintTokensThroughMarketplaceEndpointName = "mintTokensThroughMarketplace"
	mintTokensEndpointName                   = "mintTokens"
	makeOfferEndpointName                    = "makeOffer"
	acceptOfferEndpointName                  = "acceptOffer"
	cancelOfferEndpointName                  = "cancelOffer"
	startAuctionEndpointName                 = "startAuction"
	placeBidEndpointName                     = "placeBid"
	endAuctionEndpointName                   = "endAuction"
	depositEndpointName                      = "deposit"
	withdrawEndpointName                     = "withdraw"
	withdrawCreatorRoyaltiesEndpointName     = "withdrawCreatorRoyalties"
	issueNFTEndpointName                     = "issueNonFungible"
	deployNFTTemplateEndpointName            = "deployNFTTemplateContract"
	changeOwnerEndpointName                  = "changeOwner"
	setSpecialRoleEndpointName               = "setSpecialRole"
	withdrawFromMinterEndpointName           = "withdraw"
	requestWithdrawThroughMinterEndpointName = "requestWithdraw"
)

const RoyaltiesBP = 100

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
	config               config.BlockchainConfig
	noFeeOnMintContracts map[string]bool
}

func NewTxFormatter(cfg config.BlockchainConfig) TxFormatter {
	noFeeOnMintContracts := make(map[string]bool)
	for _, contract := range cfg.NoFeeOnMintContracts {
		noFeeOnMintContracts[contract] = true
	}
	return TxFormatter{
		config:               cfg,
		noFeeOnMintContracts: noFeeOnMintContracts,
	}
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

func (f *TxFormatter) NewBuyNftTxTemplate(senderAddr string, tokenId string, nonce uint64, signature []byte, price string) Transaction {
	txData := buyNftEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())
	//+ "@" + hex.EncodeToString(signature)

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

func (f *TxFormatter) MakeOfferTxTemplate(senderAddr string, tokenId string, nonce uint64, amount float64, expire uint64) Transaction {
	txData := makeOfferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(amount).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(expire)).Bytes())

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

	txData := acceptOfferEndpointName +
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
	txData := cancelOfferEndpointName +
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

func (f *TxFormatter) StartAuctionTxTemplate(senderAddr string, tokenId string, nonce uint64, minBid float64, startTime uint64, deadline uint64) (*Transaction, error) {
	marketPlaceAddress, err := data.NewAddressFromBech32String(f.config.MarketplaceAddress)
	if err != nil {
		return nil, err
	}

	txData := ESDTNFTTransferEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(1)).Bytes()) +
		"@" + hex.EncodeToString(marketPlaceAddress.AddressBytes()) +
		"@" + hex.EncodeToString([]byte(startAuctionEndpointName)) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(minBid).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(deadline)).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(startTime)).Bytes())

	return &Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   senderAddr,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.StartAuctionGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}

func (f *TxFormatter) PlaceBidTxTemplate(senderAddr string, tokenId string, nonce uint64, payment string, bidAmount float64) Transaction {
	txData := placeBidEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes()) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(bidAmount).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     payment,
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.StartAuctionGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) EndAuctionTxTemplate(senderAddr string, tokenId string, nonce uint64) Transaction {
	txData := endAuctionEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(nonce)).Bytes())

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.StartAuctionGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) DepositTxTemplate(senderAddr string, payment string) Transaction {
	txData := depositEndpointName

	return Transaction{
		Nonce:     0,
		Value:     payment,
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.DepositGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) WithdrawTxTemplate(senderAddr string, amount float64) Transaction {
	txData := withdrawEndpointName
	if amount != 0 {
		txData += "@" + hex.EncodeToString(services.GetPriceDenominated(amount).Bytes())
	}

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.WithdrawGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) WithdrawCreatorRoyaltiesTxTemplate(senderAddr string) Transaction {
	txData := withdrawCreatorRoyaltiesEndpointName

	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.MarketplaceAddress,
		SndAddr:   senderAddr,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.WithdrawGasLimit,
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
	tokenID string,
	signedMessage []byte,
) (Transaction, error) {
	endpointName := mintTokensThroughMarketplaceEndpointName
	if f.noFeeOnMintContracts[contractAddress] {
		endpointName = mintTokensEndpointName
	}
	var factor uint64 = 6
	if numberOfTokens < factor {
		factor = numberOfTokens
	}
	gasLimit := f.config.MintTokenGasLimit * (numberOfTokens/factor + 1)
	totalPrice := fmt.Sprintf("%f", mintPricePerToken*float64(numberOfTokens))
	txData := endpointName +
		"@" + hex.EncodeToString(big.NewInt(int64(numberOfTokens)).Bytes()) /*+
	"@" + hex.EncodeToString(signedMessage)*/
	col, err := storage.GetCollectionByAddr(contractAddress)
	if err != nil {
		return Transaction{}, err
	}
	tt, err := storage.GetLastNonceTokenByCollectionId(col.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tt.Nonce = 0
		} else {
			return Transaction{}, err
		}
	}
	err = collstats.AddCollectionToCheck(dtos.CollectionToCheck{CollectionAddr: contractAddress, TokenID: tokenID, Counter: int(tt.Nonce)})
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{
		Nonce:     0,
		Value:     totalPrice,
		RcvAddr:   contractAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.MintGasPrice,
		GasLimit:  gasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}

func (f *TxFormatter) NewIssueNFTTxTemplate(
	walletAddress string,
	tokenName string,
	tokenTicker string,
) Transaction {
	txData := issueNFTEndpointName +
		"@" + hex.EncodeToString([]byte(tokenName)) +
		"@" + hex.EncodeToString([]byte(tokenTicker))

	return Transaction{
		Nonce:     0,
		Value:     f.config.IssueTokenEGLDCost,
		RcvAddr:   f.config.SystemSCAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.IssueNFTGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) DeployNFTTemplateTxTemplate(
	walletAddress string,
	tokenId string,
	royalties float64,
	tokenNameBase string,
	imageBaseUrl string,
	imageExtension string,
	price float64,
	maxSupply uint64,
	saleStartTimestamp uint64,
	metadataBaseUrl string,
) Transaction {
	txData := deployNFTTemplateEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(big.NewInt(int64(royalties*RoyaltiesBP)).Bytes()) +
		"@" + hex.EncodeToString([]byte(tokenNameBase)) +
		"@" + hex.EncodeToString([]byte(imageBaseUrl)) +
		"@" + hex.EncodeToString([]byte(imageExtension)) +
		"@" + hex.EncodeToString(services.GetPriceDenominated(price).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(maxSupply)).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(int64(saleStartTimestamp)).Bytes())

	if len(metadataBaseUrl) > 1 {
		txData += "@" + hex.EncodeToString([]byte(metadataBaseUrl))
	}

	return Transaction{
		Nonce:     0,
		Value:     f.config.DeployNFTTemplateEGLDCost,
		RcvAddr:   f.config.DeployerAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.DeployNFTTemplateGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) ChangeOwnerTxTemplate(
	walletAddress string,
	contractAddress string,
) (*Transaction, error) {
	scAddress, err := data.NewAddressFromBech32String(contractAddress)
	if err != nil {
		return nil, err
	}

	txData := changeOwnerEndpointName + "@" + hex.EncodeToString(scAddress.AddressBytes())

	return &Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.DeployerAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.ChangeOwnerGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}

func (f *TxFormatter) SetSpecialRolesTxTemplate(
	walletAddress string,
	tokenId string,
	contractAddress string,
) (*Transaction, error) {
	scAddress, err := data.NewAddressFromBech32String(contractAddress)
	if err != nil {
		return nil, err
	}

	txData := setSpecialRoleEndpointName +
		"@" + hex.EncodeToString([]byte(tokenId)) +
		"@" + hex.EncodeToString(scAddress.AddressBytes()) +
		"@45534454526f6c654e4654437265617465"

	return &Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   f.config.SystemSCAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.SetSpecialRolesGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}

func (f *TxFormatter) WithdrawFromMinterTxTemplate(
	walletAddress string,
	contractAddress string,
) Transaction {
	return Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   contractAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.WithdrawFromMinterGasLimit,
		Data:      withdrawFromMinterEndpointName,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}
}

func (f *TxFormatter) RequestWithdrawThroughMinterTxTemplate(
	walletAddress string,
	contractAddress string,
) (*Transaction, error) {
	marketplaceAddress, err := data.NewAddressFromBech32String(f.config.MarketplaceAddress)
	if err != nil {
		return nil, err
	}

	txData := requestWithdrawThroughMinterEndpointName +
		"@" + hex.EncodeToString(marketplaceAddress.AddressBytes())

	return &Transaction{
		Nonce:     0,
		Value:     "0",
		RcvAddr:   contractAddress,
		SndAddr:   walletAddress,
		GasPrice:  f.config.GasPrice,
		GasLimit:  f.config.RequestWithdrawThroughMinterGasLimit,
		Data:      txData,
		Signature: "",
		ChainID:   f.config.ChainID,
		Version:   1,
		Options:   0,
	}, nil
}
