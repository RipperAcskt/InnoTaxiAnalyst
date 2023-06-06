package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/RipperAcskt/innotaxianalyst/config"
	"github.com/RipperAcskt/innotaxianalyst/internal/model"
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

	err = m.Up()
	if err != migrate.ErrNoChange && err != nil {
		return nil, fmt.Errorf("migrate up failed: %w", err)
	}

	return &ClickHouse{
			db:      db,
			migrate: m,
			cfg:     cfg},
		nil
}

func (cl *ClickHouse) WriteUser(user model.User) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cl.db.ExecContext(queryCtx, "INSERT INTO innotaxi.users (id, user_id, name, phone_number, email, raiting) VALUES($1, $2, $3, $4, $5, $6)", user.ID, user.UserID, user.Name, user.PhoneNumber, user.Email, user.Raiting)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (cl *ClickHouse) WriteDriver(driver model.Driver) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cl.db.ExecContext(queryCtx, "INSERT INTO innotaxi.drivers (id, driver_id, name, phone_number, email, raiting, taxi_type) VALUES($1, $2, $3, $4, $5, $6, $7)", driver.ID, driver.DriverID, driver.Name, driver.PhoneNumber, driver.Email, driver.Raiting, driver.TaxiType)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (cl *ClickHouse) WriteOrder(order model.Order) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cl.db.ExecContext(queryCtx, "INSERT INTO innotaxi.orders (id, order_id, user_id, driver_id, driver_name, driver_phone, driver_rating, taxi_type, from, to, date, status) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		order.ID, order.OrderID, order.UserID, order.DriverID, order.DriverName, order.DriverPhone, order.DriverRating, order.TaxiType, order.From, order.To, order.Date, order.Status)

	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}
