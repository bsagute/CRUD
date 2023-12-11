package main

import (
	"fmt"
	"time"
)

func convertTimeToEpochString(t time.Time) (string, error) {
	// Calculate epoch time in seconds
	epochTime := t.Unix()

	// Convert epoch time to string
	epochTimeString := fmt.Sprintf("%d", epochTime)

	return epochTimeString, nil
}

func main() {
	// Example time value
	currentTime := time.Now()

	// Call the convertTimeToEpochString function
	epochString, err := convertTimeToEpochString(currentTime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display the epoch time string
	fmt.Println("Epoch Time String:", epochString)
}
