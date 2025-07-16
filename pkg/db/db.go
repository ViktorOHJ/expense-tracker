package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB          *pgxpool.Pool
	ErrNotFound = fmt.Errorf("transaction not found")
)

func InitDB(parentCtx context.Context, dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	schema := `CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    is_income BOOLEAN NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

	Pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = Pool.Ping(ctx)
	if err != nil {
		Pool.Close()
		return nil, err
	}

	_, err = Pool.Exec(ctx, schema)
	if err != nil {
		Pool.Close()
		return nil, err
	}

	DB = Pool
	return Pool, nil
}
