package interaction

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

func Test_SimpleQuery(t *testing.T) {
	cfg := config.BlockchainConfig{
		ProxyUrl:            "https://devnet-gateway.elrond.com",
		ChainID:             "D",
	}

	InitBlockchainInteractor(cfg)
	bi := GetBlockchainInteractor()
	require.NotNil(t, bi)

	resp, err := bi.DoSimpleVmQuery("erd1qqqqqqqqqqqqqpgq3uvfynvpvcs8aldhuyrseuyepmp0cj7at9usgefv56", "getLeftForSale")
	require.Nil(t, err)
	require.True(t, len(resp) > 0)

	u64Hex := hex.EncodeToString(resp[0])
	leftForSale, _ := strconv.ParseUint(u64Hex, 16, 64)
	require.True(t, leftForSale != 0)
}
