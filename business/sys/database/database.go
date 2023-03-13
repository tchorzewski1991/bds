package database

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/seed.sql
	seedDoc string
)

var (
	ErrNotFound  = errors.New("entry not found")
	ErrNotUnique = errors.New("entry not unique")
)

// UniqueViolation lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const UniqueViolation = "23505"

type Config struct {
	User string
	Pass string
	Host string
	Name string
}

func Open(config Config) (*sqlx.DB, error) {
	dbURI := databaseURI(config)
	if dbURI == "" {
		return nil, errors.New("cannot construct db uri")
	}

	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Check(ctx context.Context, db *sqlx.DB) (time.Time, error) {
	// Try to ping database until no errors will be returned
	var attempts int
	for {
		attempts += 1
		err := db.Ping()
		if err == nil {
			break
		}

		select {
		case <-time.After(time.Duration(attempts) * 100 * time.Millisecond):
		case <-ctx.Done():
			return time.Time{}, ctx.Err()
		}
	}
	// We might have returned timeout in a meantime
	if ctx.Err() != nil {
		return time.Time{}, ctx.Err()
	}
	// Run simple query to make sure everything is find
	var t time.Time
	err := db.QueryRowContext(ctx, `select now()`).Scan(&t)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func Migrate(ctx context.Context, db *sqlx.DB) error {
	_, err := Check(ctx, db)
	if err != nil {
		return fmt.Errorf("checking db status failed: %w", err)
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("constructing darwin driver failed: %w", err)
	}

	return darwin.New(driver, darwin.ParseMigrations(schemaDoc)).Migrate()
}

func Seed(ctx context.Context, db *sqlx.DB) error {
	_, err := Check(ctx, db)
	if err != nil {
		return fmt.Errorf("checking db status failed: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("starting db transaction failed: %w", err)
	}

	_, err = tx.Exec(seedDoc)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

type Schema struct {
	Entries []Entry
}

type Entry struct {
	TableName  string `db:"table_name"`
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
}

func GenerateSchema(ctx context.Context, db *sqlx.DB) (Schema, error) {
	_, err := Check(ctx, db)
	if err != nil {
		return Schema{}, fmt.Errorf("checking db status failed: %w", err)
	}

	q := `select table_name, column_name, data_type from information_schema.columns where table_name = :table_name`

	args := struct {
		TableName string `db:"table_name"`
	}{TableName: "users"}

	rows, err := db.NamedQueryContext(ctx, q, args)
	if err != nil {
		return Schema{}, fmt.Errorf("executing schema information query failed: %w", err)
	}
	defer rows.Close()

	var es []Entry
	for rows.Next() {
		var e Entry
		err = rows.StructScan(&e)
		if err != nil {
			return Schema{}, fmt.Errorf("scanning table failed: %w", err)
		}
		es = append(es, e)
	}

	return Schema{
		Entries: es,
	}, nil
}

// private

func databaseURI(c Config) string {
	q := make(url.Values)
	q.Set("sslmode", "disable")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Pass),
		Host:     c.Host,
		Path:     c.Name,
		RawQuery: q.Encode(),
	}

	return u.String()
}
