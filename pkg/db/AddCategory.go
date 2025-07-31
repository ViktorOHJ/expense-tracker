package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) AddCategory(parentContext context.Context, userID int, c *models.Category) (models.Category, error) {
	query := `INSERT INTO categories (name, description, user_id)
	          VALUES ($1, $2, $3) RETURNING *`

	ctx, cancel := context.WithTimeout(parentContext, 5*time.Second)
	defer cancel()

	category := models.Category{}
	err := db.pool.QueryRow(ctx, query, c.Name, c.Description, userID).
		Scan(&category.ID, &category.Name, &category.Description, &category.UserID)

	if err != nil {
		log.Printf("failed to insert category: %v", err)
		return models.Category{}, fmt.Errorf("failed to insert category: %v", err)
	}
	return category, nil
}
