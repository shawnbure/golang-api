package process

import (
	"encoding/hex"
	data2 "github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"math/big"
	"strconv"
)

func decodeU64FromTopic(bytes []byte) uint64 {
	u64Hex := hex.EncodeToString(bytes)
	timestamp, _ := strconv.ParseUint(u64Hex, 16, 64)
	return timestamp
}

func decodeStringFromTopic(bytes []byte) string {
	return string(bytes)
}

func decodeAddressFromTopic(bytes []byte) string {
	address := data2.NewAddressFromBytes(bytes)
	return address.AddressAsBech32String()
}

func decodeBigUintFromTopic(bytes []byte) string {
	bigUint := big.NewInt(0).SetBytes(bytes)
	bigUintBytes := bigUint.Bytes()
	return hex.EncodeToString(bigUintBytes)
}

func decodeTxHashFromTopic(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
