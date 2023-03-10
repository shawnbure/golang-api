package config

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
)

type GeneralConfig struct {
	ConnectorApi       ConnectorApiConfig
	Blockchain         BlockchainConfig
	Database           DatabaseConfig
	Auth               AuthConfig
	Cache              CacheConfig
	Swagger            SwaggerConfig
	Bot                BotConfig
	Monitor            MonitorConfig
	CDN                CDNConfig
	ExternalCredential ExternalCredentialConfig
	Proxy              ProxyConfig
	CarbonSetting      CarbonSettingConfig
}

type ConnectorApiConfig struct {
	Address     string
	Username    string
	Password    string
	Addresses   []string
	Identifiers []string
}

type BlockchainConfig struct {
	GasPrice                             uint64
	ProxyUrl                             string
	ProxyUrlSec                          string
	ApiUrl                               string
	ApiUrlSec                            string
	CollectionAPIDelay                   uint64
	ChainID                              string
	PemPath                              string
	MarketplaceAddress                   string
	DeployerAddress                      string
	SystemSCAddress                      string
	StakingAddress                       string
	ListNftGasLimit                      uint64
	BuyNftGasLimit                       uint64
	WithdrawNftGasLimit                  uint64
	MintTokenGasLimit                    uint64
	MintGasPrice                         uint64
	MakeOfferGasLimit                    uint64
	AcceptOfferGasLimit                  uint64
	CancelOfferGasLimit                  uint64
	StartAuctionGasLimit                 uint64
	PlaceBidGasLimit                     uint64
	EndAuctionGasLimit                   uint64
	DepositGasLimit                      uint64
	WithdrawGasLimit                     uint64
	WithdrawCreatorRoyaltiesGasLimit     uint64
	IssueNFTGasLimit                     uint64
	DeployNFTTemplateGasLimit            uint64
	StakeNFTTemplateGasLimit             uint64
	ChangeOwnerGasLimit                  uint64
	SetSpecialRolesGasLimit              uint64
	IssueTokenEGLDCost                   string
	DeployNFTTemplateEGLDCost            string
	StakeNFTEGLDCost                     string
	WithdrawFromMinterGasLimit           uint64
	RequestWithdrawThroughMinterGasLimit uint64
	UpdateSaleStartGasLimit              uint64
	NoFeeOnMintContracts                 []string
}

type DatabaseConfig struct {
	Dialect        string
	Host           string
	Port           uint16
	DbName         string
	User           string
	Password       string
	SslMode        string
	ConnectionName string
	MaxOpenConns   int
	MaxIdleConns   int
	ShouldMigrate  bool
}

type AuthConfig struct {
	JwtSecret     string
	JwtIssuer     string
	JwtKeySeedHex string
	JwtExpiryMins int
}

type CacheConfig struct {
	ReadUrl  string
	WriteUrl string
}

type SwaggerConfig struct {
	LocalDocRoute string
	Enabled       bool
}

type BotConfig struct {
	Token  string
	RecID  string
	Enable bool
}

type MonitorConfig struct {
	ObserverMonitorEnable bool
}

type CDNConfig struct {
	Name       string
	ProjectID  string
	BucketName string
	UploadPath string
	// ApiKey    string
	// ApiSecret string
	Selector string
	BaseUrl  string
	RootDir  string
}
type ProxyConfig struct {
	List []string
}
type ExternalCredentialConfig struct {
	DreamshipAPIKey string
}

type CarbonSettingConfig struct {
	StaticAddress string
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
