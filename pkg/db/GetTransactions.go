package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func GetTransactions(ctx context.Context, transactions *[]models.Transaction) error {
	if DB == nil {
		return errors.New("DB = nil")
	}

	query := `SELECT id, is_income, amount, category_id, note, created_at FROM transactions ORDER BY created_at DESC`
	rows, err := DB.Query(ctx, query)
	if err != nil {
		log.Printf("failed to retrieve transactions: %v", err)
		return fmt.Errorf("failed to retrieve transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt); err != nil {
			log.Printf("failed to scan transaction: %v", err)
			return fmt.Errorf("failed to scan transaction: %v", err)
		}
		*transactions = append(*transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during row iteration: %v", err)
		return fmt.Errorf("error during row iteration: %v", err)
	}

	return nil
}
