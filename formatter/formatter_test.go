package formatter

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_DecodeAddress(t *testing.T) {
	bech32address := "erd1qqqqqqqqqqqqqpgq78y09lw93f3udvsplshdv2vk957l5vl70n4splrad2"
	hexAddress := ""
	_ = bech32address

	//Code here

	expectedHexAddress := "00000000000000000500f1c8f2fdc58a63c6b201fc2ed629962d3dfa33fe7ceb"
	/*
		obtinuta prin
			> erdpy wallet bech32 --decode erd1qqqqqqqqqqqqqpgq78y09lw93f3udvsplshdv2vk957l5vl70n4splrad2
			00000000000000000500f1c8f2fdc58a63c6b201fc2ed629962d3dfa33fe7ceb
	*/
	require.Equal(t, hexAddress, expectedHexAddress)
}

func TestTxFormatter_NewListNftTxTemplate(t *testing.T) {
	formatter := NewTxFormatter(defaultConfig())

	tx, err := formatter.NewListNftTxTemplate(
			"erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh",
			"LKMEX-85ea13",
			2,
			"4096",
		)

	require.Nil(t, err)
	require.Equal(t, string(tx.Data), "ESDTNFTTransfer@4C4B4D45582D383565613133@02@01@000000000000000005008D8E525546959427D05CA3172B611065D92BF3535979@7075744E6674466F7253616C65@1000")
}

func TestTxFormatter_NewWithdrawNftTxTemplate(t *testing.T) {
	formatter := NewTxFormatter(defaultConfig())

	tx := formatter.NewWithdrawNftTxTemplate(
		"erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh",
		"LKMEX-85ea13",
		2,
	)

	require.True(t, strings.EqualFold(string(tx.Data), "withdrawNft@4C4B4D45582D383565613133@02"))
}

func TestTxFormatter_NewBuyNftTxTemplate(t *testing.T) {
	formatter := NewTxFormatter(defaultConfig())

	tx := formatter.NewBuyNftTxTemplate(
		"erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh",
		"LKMEX-85ea13",
		2,
		"4096",
	)

	require.True(t, strings.EqualFold(string(tx.Data), "buyNft@4C4B4D45582D383565613133@02"))
}

func defaultConfig() config.BlockchainConfig {
	return config.BlockchainConfig{
		GasPrice:            1_000_000_000,
		ProxyUrl:            "https://devnet-gateway.elrond.com",
		ChainID:             "D",
		PemPath:             "./config/owner.pem",
		MarketplaceAddress:  "erd1qqqqqqqqqqqqqpgq3k89y42xjk2z05zu5vtjkcgsvhvjhu6nt9usruf2td",
		ListNftGasLimit:     20_000_000,
		BuyNftGasLimit:      15_000_000,
		WithdrawNftGasLimit: 15_000_000,
	}
}
