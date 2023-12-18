// save_customer_documents_test.go
package logic

import (
	"net/http"
	"testing"

	"digi-data-ingestion-client/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockMongoCollection is a mock implementation of the mongo.Collection interface for testing.
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) UpdateOne(ctx *gin.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
	args := m.Called(append([]interface{}{ctx, filter}, opts...)...)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) InsertOne(ctx *gin.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// MongoCollection is an interface representing the necessary methods from mongo.Collection.
type MongoCollection interface {
	UpdateOne(ctx *gin.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	FindOne(ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult
	InsertOne(ctx *gin.Context, document interface{}) (*mongo.InsertOneResult, error)
}

var mockFindDocData = func(collection *MockMongoCollection, ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
	result := &mongo.SingleResult{} // Create a mock SingleResult
	collection.On("FindOne", ctx, filter, opts).Return(result)
	return result
}

var mockInsertOne = func(collection MongoCollection, ctx *gin.Context, document interface{}) (*mongo.InsertOneResult, error) {
	result := &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()} // Create a mock InsertOneResult
	collection.(*MockMongoCollection).On("InsertOne", ctx, document).Return(result, nil)
	return result, nil
}

func TestSaveCustomerDocuments(t *testing.T) {
	// Mock the mongo functions
	// mongoCollection := new(MockMongoCollection)

	_ = func(collection MongoCollection, ctx *gin.Context, filter interface{}, update interface{}) error {
		return nil
	}
	_ = func(collection MongoCollection, ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
		return nil
	}
	_ = func(collection MongoCollection, ctx *gin.Context, document interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}

	// Test Case 1: Save documents successfully
	ctx := &gin.Context{}
	data := []models.DocumentData{}
	userRef := "testUserRef"
	customerID := primitive.NewObjectID()
	deviceID := primitive.NewObjectID()
	androidID := "testAndroidID"

	response := SaveCustomerDocuments(ctx, data, userRef, customerID, deviceID, androidID)

	assert.True(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)

	// Test Case 2: Duplicate documents
	_ = func(collection MongoCollection, ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
		return nil
	}

	// Call the function again with duplicate documents
	response = SaveCustomerDocuments(ctx, data, userRef, customerID, deviceID, androidID)

	assert.False(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)
	assert.Contains(t, response.Message, "Duplicate documents found.")

	// Test Case 3: Empty documents
	_ = func(collection MongoCollection, ctx *gin.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
		return nil
	}

	// Call the function with empty documents
	response = SaveCustomerDocuments(ctx, []models.DocumentData{}, userRef, customerID, deviceID, androidID)

	assert.False(t, response.Success)
	assert.Equal(t, http.StatusOK, response.ResponseCode)
	assert.Contains(t, response.Message, "No docs found to insert.")

	// Test Case 4: Error during document update
	_ = func(collection MongoCollection, ctx *gin.Context, filter interface{}, update interface{}) error {
		return errors.New("Simulated update error")
	}

	// Call the function with normal documents
	response = SaveCustomerDocuments(ctx, data, userRef, customerID, deviceID, androidID)

	assert.False(t, response.Success)
	assert.Equal(t, http.StatusInternalServerError, response.ResponseCode)
	assert.Contains(t, response.Message, "failed while updating document details in DB.")
}
