package db

import (
	"context"
	"fmt"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	AddCategory(context.Context, int, *models.Category) (models.Category, error)
	AddTransaction(context.Context, int, *models.Transaction) (models.Transaction, error)
	CheckCategory(context.Context, int, int) (bool, error) // userID, categoryID
	GetTransactions(context.Context, int, *bool, *int, *time.Time, *time.Time, int, int) ([]*models.Transaction, error)
	GetSummary(context.Context, int, time.Time, time.Time) (models.Summary, error)
	DeleteTransaction(context.Context, int, int) error                        // userID, transactionID
	GetTransactionByID(context.Context, int, int) (models.Transaction, error) // userID, transactionID
	CreateUser(context.Context, *models.User) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	GetUserByID(context.Context, int) (models.User, error)
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

	schema := `
-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица категорий с привязкой к пользователю
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, user_id) -- Уникальность имени категории в рамках пользователя
);

-- Таблица транзакций с привязкой к пользователю
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    is_income BOOLEAN NOT NULL,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_is_income ON transactions(is_income);
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Составные индексы для частых запросов
CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_income ON transactions(user_id, is_income);
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
