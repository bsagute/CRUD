// save_customer_documents_test.go
package logic

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"digi-data-ingestion-client/models"

	mongoD "go.mongodb.org/mongo-driver/mongo"
)

// Mock the mongo functions for testing
var mockUpdateDocData = func(collection *mongo.Collection, ctx *gin.Context, filter interface{}, update interface{}) error {
	// Mock implementation
	return nil
}

var mockFindDocData = func(collection *mongo.Collection, ctx *gin.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	// Mock implementation
	// You can customize the mock behavior based on your test cases
	return &mongo.SingleResult{}
}

func TestSaveCustomerDocuments(t *testing.T) {
	// Mock the mongo functions
	UpdateDocData = mockUpdateDocData
	FindDocData = mockFindDocData

	// Test Case 1: Save documents successfully
	ctx := &gin.Context{} // Mock Gin context
	data := []models.DocumentData{
		// Provide your test data
	}
	userRef := "testUserRef"
	customerID := primitive.NewObjectID()
	deviceID := primitive.NewObjectID()
	androidID := "testAndroidID"

	response := SaveCustomerDocuments(ctx, data, userRef, customerID, deviceID, androidID)

	assert.True(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)

	// Test Case 2: Duplicate documents
	// Set up mock function to simulate finding existing documents
	mockFindDocData = func(collection *mongoD.Collection, ctx *gin.Context, filter interface{}, opts ...*options.FindOneOptions) *mongoD.SingleResult {
		return &mongoD.SingleResult{} // Simulate finding existing documents
	}
	FindDocData = mockFindDocData

	// Call the function again with duplicate documents
	response = SaveCustomerDocuments(ctx, data, userRef, customerID, deviceID, androidID)

	assert.False(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)
	assert.Contains(t, response.Message, "Duplicate documents found.")

	// Test Case 3: Empty documents
	// Modify the mock function to simulate finding no existing documents
	mockFindDocData = func(collection *mongoD.Collection, ctx *gin.Context, filter interface{}, opts ...*options.FindOneOptions) *mongoD.SingleResult {
		return nil // Simulate not finding any existing documents
	}
	FindDocData = mockFindDocData

	// Call the function with empty documents
	response = SaveCustomerDocuments(ctx, []models.DocumentData{}, userRef, customerID, deviceID, androidID)

	assert.False(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)            // Update expected status code
	assert.Contains(t, response.Message, "No docs found to insert.") // Update expected error message
}
