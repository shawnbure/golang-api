package config

import "github.com/ElrondNetwork/elrond-go-core/core"

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	Blockchain   BlockchainConfig
}

type ConnectorApiConfig struct {
	Port      string
	Username  string
	Password  string
	Addresses []string
}

type BlockchainConfig struct {
	GasPrice uint64
	ProxyUrl string
	ChainID  string
	PemPath  string
}

func LoadConfig(filePath string) (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
