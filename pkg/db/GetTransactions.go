package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func GetTransactions(ctx context.Context, txType *bool, category_id *int, from, to *time.Time) ([]*models.Transaction, error) {

	query := `SELECT * FROM transactions WHERE 1=1`
	args := []interface{}{}
	i := 1
	if txType != nil {
		query += fmt.Sprintf(` AND is_income = $%d`, i)
		args = append(args, *txType)
		i++
	}
	if category_id != nil {
		query += fmt.Sprintf(` AND category_id = $%d`, i)
		args = append(args, *category_id)
		i++
	}
	if from != nil {
		query += fmt.Sprintf(` AND created_at >= $%d`, i)
		args = append(args, *from)
		i++
	}
	if to != nil {
		query += fmt.Sprintf(` AND created_at <= $%d`, i)
		args = append(args, *to)
		i++
	}
	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to retrieve transactions: %v", err)
		return []*models.Transaction{}, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt)
		if err != nil {
			log.Printf("failed to scan transaction: %v", err)
			return []*models.Transaction{}, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error iterating over rows: %v", err)
		return []*models.Transaction{}, err
	}
	return transactions, nil
}
