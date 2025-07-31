package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) AddTransaction(parentCtx context.Context, userID int, t *models.Transaction) (models.Transaction, error) {
	query := `INSERT INTO transactions (is_income, amount, category_id, user_id, note)
	          VALUES ($1, $2, $3, $4, $5) RETURNING *`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	transaction := models.Transaction{}
	err := db.pool.QueryRow(ctx, query, t.IsIncome, t.Amount, t.CategoryID, userID, t.Note).
		Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount,
			&transaction.CategoryID, &transaction.UserID, &transaction.Note, &transaction.CreatedAt)

	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return models.Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return transaction, nil
}

func (db *PostgresDB) CheckCategory(parentCtx context.Context, userID int, categoryID int) (exists bool, err error) {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE id=$1 AND user_id=$2)`

	err = db.pool.QueryRow(ctx, query, categoryID, userID).Scan(&exists)
	if err != nil {
		log.Printf("database error during category check: %v", err)
		return false, fmt.Errorf("database error during category check: %v", err)
	}
	return exists, nil
}
