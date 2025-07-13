package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func AddTransaction(parentCtx context.Context, t *models.Transaction) (models.Transaction, error) {
	if DB == nil {
		log.Println("db = nil")
		return models.Transaction{}, errors.New("DB = nil")
	}
	query := `INSERT INTO transactions (is_income, amount, category_id, note) VALUES ($1, $2, $3, $4)
RETURNING *`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	transaction := models.Transaction{}
	err := DB.QueryRow(ctx, query, t.IsIncome, t.Amount, t.CategoryID, t.Note).Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return models.Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return transaction, nil
}

func CheckCategory(parentCtx context.Context, catID int) (exists bool, err error) {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE id=$1)`

	err = DB.QueryRow(ctx, query, catID).Scan(&exists)
	if err != nil {
		log.Printf("database error during category check: %v", err)
		return false, fmt.Errorf("database error during category check: %v", err)
	}
	if !exists {
		return false, nil
	}
	return exists, nil
}
