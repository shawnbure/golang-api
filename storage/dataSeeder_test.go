package storage

import (
	"testing"

	"github.com/erdsea/erdsea-api/config"
)

var cfg = config.DatabaseConfig{
	Dialect:       "postgres",
	Host:          "localhost",
	Port:          5432,
	DbName:        "erdsea_db_test",
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
