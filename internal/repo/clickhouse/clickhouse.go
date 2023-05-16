package clickhouse

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/golang-migrate/migrate/v4"
	migrateClickHouse "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type ClickHouse struct {
	db      *sql.DB
	migrate *migrate.Migrate
	cfg     *config.Config
}

func New(cfg *config.Config) (*ClickHouse, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cfg.CLICKHOUSE_DB_HOST},
		Auth: clickhouse.Auth{
			Database: cfg.CLICKHOUSE_DB_NAME,
			Username: cfg.CLICKHOUSE_DB_USERNAME,
			Password: cfg.CLICKHOUSE_DB_PASSWORD,
		},
	})

	driver, err := migrateClickHouse.WithInstance(db, &migrateClickHouse.Config{})
	if err != nil {
		return nil, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MIGRATE_PATH, "clickhouse", driver)
	if err != nil {
		return nil, fmt.Errorf("new with database instance failed: %w", err)
	}

	return &ClickHouse{
			db:      db,
			migrate: m,
			cfg:     cfg},
		nil
}
