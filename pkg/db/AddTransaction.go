package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) AddTransaction(parentCtx context.Context, t *models.Transaction) (models.Transaction, error) {
	query := `INSERT INTO transactions (is_income, amount, category_id, note) VALUES ($1, $2, $3, $4)
RETURNING *`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	transaction := models.Transaction{}
	err := db.pool.QueryRow(ctx, query, t.IsIncome, t.Amount, t.CategoryID, t.Note).Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return models.Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return transaction, nil
}

func (db *PostgresDB) CheckCategory(parentCtx context.Context, id int) (exists bool, err error) {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE id=$1)`

	err = db.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		log.Printf("database error during category check: %v", err)
		return false, fmt.Errorf("database error during category check: %v", err)
	}
	return exists, nil
}
