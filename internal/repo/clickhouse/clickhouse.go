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

	_, err := cl.db.ExecContext(queryCtx, "INSERT INTO innotaxi.users (id, user_id, name, phone_number, email, rating) VALUES($1, $2, $3, $4, $5, $6)", user.ID, user.UserID, user.Name, user.PhoneNumber, user.Email, user.Raiting)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (cl *ClickHouse) WriteDriver(driver model.Driver) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cl.db.ExecContext(queryCtx, "INSERT INTO innotaxi.drivers (id, driver_id, name, phone_number, email, rating, taxi_type) VALUES($1, $2, $3, $4, $5, $6, $7)", driver.ID, driver.DriverID, driver.Name, driver.PhoneNumber, driver.Email, driver.Raiting, driver.TaxiType)
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

func (cl *ClickHouse) SetRatingUser(ctx context.Context, r model.Rating) (float64, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := cl.db.QueryRowContext(queryCtx, "SELECT rating, num_of_marks FROM innotaxi.users WHERE user_id = $1", r.ID)
	if row.Err() != nil {
		return 0, fmt.Errorf("query row context failed: %w", row.Err())
	}

	var rating float64
	var numOfMarks int64
	err := row.Scan(&rating, &numOfMarks)
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}

	numOfMarks++
	rating += float64(r.Rating)

	_, err = cl.db.ExecContext(queryCtx, "ALTER TABLE innotaxi.users UPDATE rating = $1, num_of_marks = $2 WHERE user_id = $3", rating, numOfMarks, r.ID)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}
	return rating / float64(numOfMarks), nil
}

func (cl *ClickHouse) SetRatingDriver(ctx context.Context, r model.Rating) (float64, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := cl.db.QueryRowContext(queryCtx, "SELECT rating, num_of_marks FROM innotaxi.drivers WHERE driver_id = $1", r.ID)
	if row.Err() != nil {
		return 0, fmt.Errorf("query row context failed: %w", row.Err())
	}

	var rating float64
	var numOfMarks int64
	err := row.Scan(&rating, &numOfMarks)
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}

	numOfMarks++
	rating += float64(r.Rating)

	_, err = cl.db.ExecContext(queryCtx, "ALTER TABLE innotaxi.drivers UPDATE rating = $1, num_of_marks = $2 WHERE driver_id = $3", rating, numOfMarks, r.ID)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}
	return rating / float64(numOfMarks), nil
}

func (cl *ClickHouse) GetRating(ctx context.Context, db string) ([]model.Rating, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cl.db.QueryContext(queryCtx, "SELECT id, rating, FROM innotaxi.$1", db)
	if err != nil {
		return nil, fmt.Errorf("query row context failed: %w", err)
	}

	ratings := make([]model.Rating, 5, 10)
	for rows.Next() {
		var rating model.Rating
		err := rows.Scan(&rating.ID, rating.Rating)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
	}

	return ratings, nil
}
