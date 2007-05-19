package services

import (
	"fmt"
)

type ListAssetArgs struct {
	OwnerAddress     string
	TokenId          string
	Nonce            uint64
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

func (args *ListAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"RoyaltiesPercent = %d\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.RoyaltiesPercent,
		args.Timestamp,
		args.TxHash)
}

func (args *BuyAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"BuyerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.BuyerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.Timestamp,
		args.TxHash)
}

func (args *WithdrawAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.Timestamp,
		args.TxHash)
}
