package process

import (
	"encoding/hex"
	"math/big"
	"strconv"

	erdData "github.com/ElrondNetwork/elrond-sdk-erdgo/data"
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
	address := erdData.NewAddressFromBytes(bytes)
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

func decodeHexStringOrEmptyWhenZeroFromTopic(bytes []byte) string {
	if allZero(bytes) {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}