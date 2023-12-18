package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Assuming you already have a PostgreSQL connection, replace the connection string with your own.
	connectionString := "user=username dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Assuming you have a single map with string keys and various types of values
	data := map[string]interface{}{
		"column_string1": "value1",
		"column_string2": "value2",
		"column_string3": "value3",
		"column_bool":    true,
		"column_float":   3.14,
	}

	// Construct the INSERT query
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
	}

	query := fmt.Sprintf("INSERT INTO your_table (%s) VALUES (%s);",
		strings.Join(columns, ", "), createPlaceholders(len(values)))

	// Execute the query
	_, err = db.Exec(query, values...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Insert successful!")
}

// Helper function to create placeholders for the VALUES part of the query
func createPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}
