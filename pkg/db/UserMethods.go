package db

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5"
)

func (db *PostgresDB) CreateUser(parentCtx context.Context, user *models.User) (models.User, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	var newUser models.User
	err := db.pool.QueryRow(ctx, query, user.Email, user.Password).
		Scan(&newUser.ID, &newUser.Email, &newUser.CreatedAt)

	if err != nil {
		log.Printf("failed to create user: %v", err)
		return models.User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return newUser, nil
}

func (db *PostgresDB) GetUserByEmail(parentCtx context.Context, email string) (models.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	var user models.User
	err := db.pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (db *PostgresDB) GetUserByID(parentCtx context.Context, id int) (models.User, error) {
	query := `SELECT id, email, created_at FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	var user models.User
	err := db.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
