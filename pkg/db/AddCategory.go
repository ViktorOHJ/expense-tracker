package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (db *PostgresDB) AddCategory(parentContext context.Context, c *models.Category) (models.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING *`
	ctx, cancel := context.WithTimeout(parentContext, 5*time.Second)
	defer cancel()
	category := models.Category{}
	err := db.pool.QueryRow(ctx, query, c.Name, c.Description).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return models.Category{}, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return category, nil
}
