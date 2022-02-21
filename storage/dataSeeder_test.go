package storage

import (
	"testing"

	"github.com/ENFT-DAO/youbei-api/config"
)

var cfg = config.DatabaseConfig{
	Dialect:       "postgres",
	Host:          "localhost",
	Port:          5432,
	DbName:        "youbei_dev",
	User:          "postgres",
	Password:      "boop",
	SslMode:       "disable",
	MaxOpenConns:  50,
	MaxIdleConns:  10,
	ShouldMigrate: true,
}

func TestSeedDatabase(t *testing.T) {
	SeedDatabase(cfg)
}
