package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

func generateRecord() []string {
	// Generate random four-digit value for MetricValue
	rand.Seed(time.Now().UnixNano())
	metricValue := float64(rand.Intn(9000) + 1000)

	// Generate new GUIDs for EntityID, MetricId
	entityID := uuid.New().String()
	metricID := uuid.New().String()

	// Generate timestamp until seconds
	timestamp := time.Now().Format("2006-01-02T15:04:05Z")

	// Prepare record for CSV
	record := []string{
		entityID, // EntityID
		strconv.FormatFloat(metricValue, 'f', -1, 64), // MetricValue
		metricID,  // MetricId
		timestamp, // Timestamp
	}

	return record
}

func writeRecords(filename string, numRows int) {
	// Create the CSV file with appended count in the filename
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"EntityID", "MetricValue", "MetricId", "Timestamp"}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
		return
	}

	// Create a channel to pass row data
	ch := make(chan []string)

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start goroutines to generate records
	for i := 0; i < numRows; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- generateRecord()
		}()
	}

	// Receive row data from the channel and write to CSV
	go func() {
		for record := range ch {
			if err := writer.Write(record); err != nil {
				fmt.Printf("Error writing record to CSV: %v\n", err)
				return
			}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the channel
	close(ch)
}

func main() {
	startTime := time.Now()

	// Ask user for the number of rows
	fmt.Print("Enter the number of rows for the CSV file: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// Convert input to an integer
	numRows, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter a valid number.")
		return
	}

	filename := fmt.Sprintf("records_%d.csv", numRows)
	writeRecords(filename, numRows)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Printf("CSV file '%s' with %d records generated successfully in %.2f seconds!\n", filename, numRows, elapsed.Seconds())
}
