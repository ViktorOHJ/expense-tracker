package db

import (
	"context"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) GetSummary(parentCtx context.Context, userID int, from, to time.Time) (models.Summary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN is_income THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN NOT is_income THEN amount ELSE 0 END), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN is_income THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN NOT is_income THEN amount ELSE 0 END), 0) AS balance
		FROM transactions
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	var summary models.Summary
	err := db.pool.QueryRow(ctx, query, userID, from, to).
		Scan(&summary.TotalIncome, &summary.TotalExpense, &summary.Balance)
	if err != nil {
		log.Printf("failed to retrieve summary: %v", err)
		return models.Summary{}, err
	}
	log.Printf("Summary retrieved successfully for user %d: %+v", userID, summary)
	return summary, nil
}
