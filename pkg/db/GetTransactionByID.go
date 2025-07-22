package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5"
)

func (db *PostgresDB) GetTransactionByID(parentCtx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction
	query := `SELECT * FROM transactions WHERE id=$1`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()
	row := db.pool.QueryRow(ctx, query, id)

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
