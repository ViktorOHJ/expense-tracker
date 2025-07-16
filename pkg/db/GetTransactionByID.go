package db

import (
	"context"
	"fmt"
	"log"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5"
)

func GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction
	query := `SELECT * FROM transactions WHERE id=$1`
	row := DB.QueryRow(ctx, query, id)

	err := row.Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("transaction with id %d not found", id)
			return models.Transaction{}, ErrNotFound // No transaction found
		}
		log.Printf("failed to scan: %v", err)
		return models.Transaction{}, fmt.Errorf("failed to retrieve transaction: %v", err)
	}
	return transaction, nil
}
