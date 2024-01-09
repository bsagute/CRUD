package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

const (
	tableName = "your_table"
)

var db *sql.DB

func main() {
	// Replace these connection parameters with your PostgreSQL credentials
	connStr := "user=username dbname=mydb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the correct database is selected
	_, err = db.Exec("USE mydb")
	if err != nil {
		log.Fatal(err)
	}

	// Check if the table exists, create it if not
	err = createTableIfNotExists(tableName)
	if err != nil {
		log.Fatal(err)
	}

	// Your input values
	hashKey := "your_hash_key"
	details := map[string]interface{}{
		"column1": "value1",
		"column2": 42,
		// Add other key-value pairs as needed
	}

	// Insert the record
	err = InsertRecordInPostgres(hashKey, details)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully.")
}

// InsertRecordInPostgres inserts a record into the specified table
func InsertRecordInPostgres(hashKey string, details map[string]interface{}) error {
	// Prepare the SQL statement
	query := "INSERT INTO " + tableName + "(hash_key, "

	// Prepare placeholders for the columns and values
	var columns, placeholders []string
	for column, value := range details {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
	}

	query += fmt.Sprintf("%s) VALUES ('%s', %s)",
		strings.Join(columns, ", "),
		hashKey,
		strings.Join(placeholders, ", "),
	)

	// Execute the SQL statement
	values := make([]interface{}, 0, len(details))
	for _, value := range details {
		values = append(values, value)
	}

	_, err := db.Exec(query, values...)
	if err != nil {
		log.Println("Error inserting record:", err)
	}
	return err
}

// createTableIfNotExists checks if the specified table exists, creates it if not
func createTableIfNotExists(tableName string) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (hash_key TEXT, column1 TEXT, column2 INTEGER)")
	if err != nil {
		log.Println("Error creating table:", err)
	}
	return err
}
