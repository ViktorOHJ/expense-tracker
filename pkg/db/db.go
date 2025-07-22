package db

import (
	"context"
	"fmt"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	AddCategory(context.Context, *models.Category) (models.Category, error)
	AddTransaction(context.Context, *models.Transaction) (models.Transaction, error)
	CheckCategory(context.Context, int) (bool, error)
	GetTransactions(parentCtx context.Context, txType *bool, category_id *int, from, to *time.Time, limit, offset int) ([]*models.Transaction, error)
	GetSummary(parentCtx context.Context, from, to time.Time) (models.Summary, error)
	DeleteTransaction(context.Context, int) error
	GetTransactionByID(context.Context, int) (models.Transaction, error)
}

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pool: pool}
}

var ErrNotFound = fmt.Errorf("transaction not found")

func InitDB(parentCtx context.Context, dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	schema := `CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    is_income BOOLEAN NOT NULL,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    category_id INTEGER REFERENCES categories(id) ON DELETE RESTRICT,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_is_income ON transactions(is_income);
`

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	_, err = pool.Exec(ctx, schema)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
