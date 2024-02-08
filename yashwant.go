package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection URI
	uri := "mongodb://<username>:<password>@<host>:<port>/?authSource=admin"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Access the database and collections
	db := client.Database("authentication")
	authProfileCollection := db.Collection("authProfile")
	authDBCollection := db.Collection("authDB")

	// Define the filter to match documents where 'created' field does not exist
	filter := bson.M{"created": bson.M{"$exists": false}}

	// Find documents in 'authProfile' collection
	cursor, err := authProfileCollection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Error finding documents in authProfile collection: %v", err)
	}
	defer cursor.Close(ctx)

	// Slice to store customerIds
	var customerIDs []string

	// Iterate over the cursor and append customerId to the slice
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding result: %v", err)
		}
		customerID := result["customerId"].(string)
		customerIDs = append(customerIDs, customerID)
	}

	// Find documents in 'authDB' collection
	cursor, err = authDBCollection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Error finding documents in authDB collection: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and remove customerId from the slice if found in authDB
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding result: %v", err)
		}
		account := result["account"].(bson.M)
		customerID := account["customerId"].(string)
		for i, id := range customerIDs {
			if id == customerID {
				customerIDs = append(customerIDs[:i], customerIDs[i+1:]...)
				break
			}
		}
	}

	// Print the list of customerIds
	fmt.Println("CustomerIds with missing 'created' field in both collections:")
	for _, id := range customerIDs {
		fmt.Println(id)
	}
}
