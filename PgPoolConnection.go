// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/your_username/your_project/db"
)

func main() {
	// Replace these connection details with your PostgreSQL server information
	connConfig := pgxpool.Config{
		ConnConfig: pgxpool.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "your_database",
			User:     "your_username",
			Password: "your_password",
		},
		MaxConnLifetime: 5 * time.Minute,
		MaxConns:        5,
		MinConns:        2,
	}

	// Create a pool
	pool, err := pgxpool.ConnectConfig(context.Background(), &connConfig)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	defer pool.Close()

	// Example insertion
	emp := db.Emp{
		Name:       "John Doe",
		CreatedAt:  time.Now(),
		IsExist:    true,
		FloatValue: 123.45,
	}

	if err := db.InsertEmployee(pool, emp); err != nil {
		log.Fatalf("Insert operation failed: %v", err)
	}

	// Rest of your code
}
