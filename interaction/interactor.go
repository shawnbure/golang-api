package interaction

import (
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/config"
)

var log = logger.GetOrCreate("interaction")

const (
	waitForTxDuration = time.Second * 2
)

type BlockchainInteractor struct {
	proxyUrl   string
	chainID    string
	gasPrice   uint64
	privateKey []byte
	publicKey  string
	proxy      blockchain.ProxyHandler
	account    data.Account
	txMut      sync.Mutex
}

func NewBlockchainInteractor(chainInfo config.BlockchainConfig) (*BlockchainInteractor, error) {
	sk, pk, err := GetKeyPairFromPem(chainInfo.PemPath)
	if err != nil {
		return nil, err
	}

	proxy := blockchain.NewElrondProxy(chainInfo.ProxyUrl, &http.Client{})

	addressHandler, err := data.NewAddressFromBech32String(pk)
	if err != nil {
		return nil, err
	}

	account, err := proxy.GetAccount(addressHandler)
	if err != nil {
		return nil, err
	}

	return &BlockchainInteractor{
		proxyUrl:   chainInfo.ProxyUrl,
		chainID:    chainInfo.ChainID,
		gasPrice:   chainInfo.GasPrice,
		account:    *account,
		privateKey: sk,
		publicKey:  pk,
		proxy:      proxy,
	}, nil
}

func (bi *BlockchainInteractor) SendTx(tx *data.Transaction) (string, error) {
	txHash, err := bi.proxy.SendTransaction(tx)
	if err != nil {
		log.Debug("failed sending transaction", "err", err.Error())
		return "", err
	}

	log.Info("current local nonce", "nonce", tx.Nonce)
	return txHash, nil
}

func (bi *BlockchainInteractor) CreateSignedTx(
	inputData []byte,
	receiver string,
	gasLimit uint64,
) (*data.Transaction, error) {
	bi.txMut.Lock()
	defer bi.txMut.Unlock()

	account, err := bi.getAccount()
	if err != nil {
		return nil, err
	}
	if account.Nonce > bi.account.Nonce {
		log.Debug("got higher nonce from proxy",
			"nonce", account.Nonce,
			"current nonce", bi.account.Nonce,
			"replacing", true,
		)
		bi.account.Nonce = account.Nonce
	}

	tx := &data.Transaction{
		Value:    "0",
		RcvAddr:  receiver,
		Data:     inputData,
		SndAddr:  account.Address,
		Nonce:    bi.account.Nonce,
		ChainID:  bi.chainID,
		GasPrice: bi.gasPrice,
		GasLimit: gasLimit,
		Version:  1,
		Options:  0,
	}

	err = erdgo.SignTransaction(tx, bi.privateKey)
	if err != nil {
		log.Debug("failed signing transaction", "err", err.Error())
		return nil, err
	}

	bi.account.Nonce++
	return tx, nil
}

func (bi *BlockchainInteractor) WaitTxExecutionCheckStatus(txHash string) (bool, error) {
	for {
		time.Sleep(waitForTxDuration)
		status, err := bi.proxy.GetTransactionStatus(txHash)
		if err != nil {
			return false, err
		}

		switch status {
		case transaction.TxStatusPending.String():
			continue
		case transaction.TxStatusSuccess.String():
			return true, nil
		case transaction.TxStatusInvalid.String(), transaction.TxStatusFail.String():
			return false, nil
		}
	}
}

func (bi *BlockchainInteractor) getAccount() (*data.Account, error) {
	addressHandler, err := data.NewAddressFromBech32String(bi.publicKey)
	if err != nil {
		return nil, err
	}

	account, err := bi.proxy.GetAccount(addressHandler)
	if err != nil {
		return nil, err
	}
	return account, nil
}
