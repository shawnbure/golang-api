package process

import (
	"encoding/hex"
	logger "log"
	"os"
	"strconv"
	"time"

	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/storage"
)

type EndAuctionChecker struct {
	MarketPlaceAddr string `json:"marketPlaceAddr"`
	ElrondAPI       string `json:"elrondAPI"`
	Logger          *logger.Logger
	Delay           time.Duration // delay between each call
}

func NewEndAuctionChecker(marketPlaceAddr string, elrondAPI string, delay uint64) (*EndAuctionChecker, error) {
	lerr := logger.New(os.Stderr, "", logger.LUTC|logger.LstdFlags|logger.Lshortfile)
	return &EndAuctionChecker{MarketPlaceAddr: marketPlaceAddr, ElrondAPI: elrondAPI, Logger: lerr, Delay: time.Duration(delay)}, nil
}

func (mpi *EndAuctionChecker) StartWorker() {
	lerr := mpi.Logger
	bi := interaction.GetBlockchainInteractor()
	for {
		tokens, err := storage.GetEndAuctionTokens()
		if err != nil {
			lerr.Println(err.Error())
			continue
		}
		for _, token := range tokens {
			nonceStr := strconv.FormatUint(token.Nonce, 10)
			nonceHex := hex.EncodeToString([]byte(nonceStr))
			tokenIdHex := hex.EncodeToString([]byte(token.TokenID))

			res, err := bi.DoVmQuery(mpi.MarketPlaceAddr, "endAuction", []string{nonceHex, tokenIdHex})
			if err != nil {
				lerr.Println(err.Error())
				continue
			}

			if len(res) == 0 {
				lerr.Println("no response from endAuction call")
			}

		}

	}

}
