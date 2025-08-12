package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) GetTransactions(parentCtx context.Context, userID int, txType *bool, category_id *int, from, to *time.Time, limit, offset int) ([]*models.Transaction, error) {

	query := `SELECT * FROM transactions WHERE user_id = $1`
	args := []interface{}{userID}
	i := 2 // Start with 2 because $1 is already used for userID

	if txType != nil {
		query += ` AND is_income = $` + strconv.Itoa(i)
		args = append(args, *txType)
		i++
	}
	if category_id != nil {
		query += ` AND category_id = $` + strconv.Itoa(i)
		args = append(args, *category_id)
		i++
	}
	if from != nil {
		query += ` AND created_at >= $` + strconv.Itoa(i)
		args = append(args, *from)
		i++
	}
	if to != nil {
		query += ` AND created_at <= $` + strconv.Itoa(i)
		args = append(args, *to)
		i++
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, i, i+1)
	args = append(args, limit, offset)

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to retrieve transactions: %v", err)
		return []*models.Transaction{}, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.IsIncome, &transaction.Amount,
			&transaction.CategoryID, &transaction.UserID, &transaction.Note, &transaction.CreatedAt)
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
