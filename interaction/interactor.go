package interaction

import (
	"net/http"
	"sync"

	"github.com/ENFT-DAO/youbei-api/config"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("interaction")

type BlockchainInteractor struct {
	proxyUrl string
	chainID  string
	proxy    blockchain.ProxyHandler
}

var (
	interactor *BlockchainInteractor
	once       sync.Once
)

func InitBlockchainInteractor(chainInfo config.BlockchainConfig) {
	once.Do(func() {
		proxy := blockchain.NewElrondProxy(chainInfo.ProxyUrl, &http.Client{})
		interactor = &BlockchainInteractor{
			proxyUrl: chainInfo.ProxyUrl,
			chainID:  chainInfo.ChainID,
			proxy:    proxy,
		}
	})
}

func (bi *BlockchainInteractor) DoVmQuery(contractAddress string, viewFuncName string, args []string) ([][]byte, error) {
	request := data.VmValueRequest{
		Address:  contractAddress,
		FuncName: viewFuncName,
		Args:     args,
	}

	response, err := bi.proxy.ExecuteVMQuery(&request)
	if err != nil {
		log.Debug("failed to execute vm query", err.Error())
		return nil, err
	}

	return response.Data.ReturnData, nil
}

func GetBlockchainInteractor() *BlockchainInteractor {
	return interactor
}
