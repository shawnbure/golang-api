package config

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
)

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	Blockchain   BlockchainConfig
	Database     DatabaseConfig
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

type DatabaseConfig struct {
	Dialect       string
	Host          string
	Port          uint16
	DbName        string
	User          string
	Password      string
	SslMode       string
	MaxOpenConns  int
	MaxIdleConns  int
	ShouldMigrate bool
}

func (d DatabaseConfig) Url() string {
	format := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return fmt.Sprintf(format, d.Host, d.Port, d.User, d.Password, d.DbName, d.SslMode)
}

func LoadConfig(filePath string) (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
