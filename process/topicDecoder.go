package process

import (
	"encoding/hex"
	data2 "github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"strconv"
)

func decodeU64FromTopic(bytes []byte) uint64 {
	timestamp, _ := strconv.ParseUint(hex.EncodeToString(bytes), 16, 64)
	return timestamp
}

func decodeStringFromTopic(bytes []byte) string {
	return string(bytes)
}

func decodeAddressFromTopic(bytes []byte) string {
	return data2.NewAddressFromBytes(bytes).AddressAsBech32String()
}
