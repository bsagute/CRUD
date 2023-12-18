package mongo

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Custom error variable for document not found
var ErrDocumentNotFound = fmt.Errorf("document not found")
var ErrNoDocuments = mongo.ErrNoDocuments

// Create collection wise variables
var (
	AuditCollection                     *mongo.Collection
	CustomerDevicesCollection           *mongo.Collection
	CustomersCollection                 *mongo.Collection
	CustomerLocationsCollection         *mongo.Collection
	CustomerDeviceMapperCollection      *mongo.Collection
	CustomerBankTransactions            *mongo.Collection
	CustomerEmis                        *mongo.Collection
	CustomerCompanyExpenses             *mongo.Collection
	CustomerOrdersData                  *mongo.Collection
	CustomerInvestmentsData             *mongo.Collection
	CustomerTaxData                     *mongo.Collection
	CustomerInsuranceTransactions       *mongo.Collection
	CustomerTradingTransactions         *mongo.Collection
	CustomerOnboardMessages             *mongo.Collection
	CustomerComplaintTransactions       *mongo.Collection
	CustomerLegalNoticeMessages         *mongo.Collection
	CustomerTransactionStatusCollection *mongo.Collection
	INITCUSTOMERMESSAGES                *mongo.Collection
	CustomerMessagesCollection          *mongo.Collection
	CustomerAppsCollection              *mongo.Collection
	AppCategoriesCollection             *mongo.Collection
)
var SaveMessagesCollectionsToCheck map[string]interface{}

// Returns a handle for a collection
func getCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	return collection
}

// Returns a handle for a collection
func getDbDumpCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(os.Getenv("MONGO_DB_DUMP")).Collection(collectionName)
	return collection
}

// Initialize collection objects
func configure(client *mongo.Client) {
	AuditCollection = getCollection(client, "audit")
	CustomerDevicesCollection = getCollection(client, "customer_devices")
	CustomerLocationsCollection = getCollection(client, "customer_locations")
	CustomersCollection = getCollection(client, "customers")
	CustomerDeviceMapperCollection = getCollection(client, "customer_device_mapper")
	CustomerBankTransactions = getCollection(client, "customer_bank_transactions")
	CustomerEmis = getCollection(client, "customer_emis")
	CustomerCompanyExpenses = getCollection(client, "customer_company_expenses")
	CustomerOrdersData = getCollection(client, "customer_orders_data")
	CustomerInvestmentsData = getCollection(client, "customer_investments_data")
	CustomerTaxData = getCollection(client, "customer_tax_data")
	CustomerInsuranceTransactions = getCollection(client, "customer_insurance_transactions")
	CustomerTradingTransactions = getCollection(client, "customer_trading_transactions")
	CustomerOnboardMessages = getCollection(client, "customer_onboard_messages")
	CustomerComplaintTransactions = getCollection(client, "customer_complaint_transactions")
	CustomerLegalNoticeMessages = getCollection(client, "customer_legal_notice_messages")
	CustomerTransactionStatusCollection = getCollection(client, "customer_transaction_status")
	CustomerMessagesCollection = getCollection(client, "customer_messages")
	CustomerAppsCollection = getCollection(client, "customer_apps")
	AppCategoriesCollection = getCollection(client, "app_categories")

	SaveMessagesCollectionsToCheck = map[string]interface{}{

		"CustomerBankTransactions":      CustomerBankTransactions,
		"CustomerEmis":                  CustomerEmis,
		"CustomerCompanyExpenses":       CustomerCompanyExpenses,
		"CustomerOrdersData":            CustomerOrdersData,
		"CustomerInvestmentsData":       CustomerInvestmentsData,
		"CustomerTaxData":               CustomerTaxData,
		"CustomerInsuranceTransactions": CustomerInsuranceTransactions,
		"CustomerTradingTransactions":   CustomerTradingTransactions,
		"CustomerOnboardMessages":       CustomerOnboardMessages,
		"CustomerComplaintTransactions": CustomerComplaintTransactions,
		"CustomerLegalNoticeMessages":   CustomerLegalNoticeMessages,
	}

}

// Initialize Db Dump collections objects
func dbDumpConfigure(client *mongo.Client) {
	INITCUSTOMERMESSAGES = getDbDumpCollection(client, "init_customer_messages")
}

// Inserts a single document in given collection
func InsertOne(collection *mongo.Collection, ctx *gin.Context, document interface{}) (primitive.ObjectID, error) {

	req, err := collection.InsertOne(ctx, document)
	if err != nil {
		panic(err)
	}
	insertedId := req.InsertedID
	return insertedId.(primitive.ObjectID), nil
}

// Updates a single document in given collection
func UpdateOne(collection *mongo.Collection, ctx *gin.Context, filter interface{}, update interface{}) error {

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}

// Find multiple documents based on given filters from given collection
func Find(collection *mongo.Collection, ctx *gin.Context, filter interface{}, opts *options.FindOptions) (*mongo.Cursor, error) {

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		panic(err)
	}
	return cursor, err
}

// Finds a single document from given collection
func FindOne(collection *mongo.Collection, ctx *gin.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {

	result := collection.FindOne(ctx, filter, opts...)
	return result
}

// Return count of documents from given collection
func Count(collection *mongo.Collection, ctx *gin.Context, filter interface{}) (int64, error) {

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}
	return count, err
}

// FindOne retrieves a document from the collection by its ID with projection.
// It takes the document ID as input, a projection filter, and a pointer to the variable where the result will be stored.
// It returns an error if the document is not found or any other error occurs.
func FindOneWithProjection(collection *mongo.Collection, ctx *gin.Context, filter bson.M, projection bson.M, doc interface{}) error {

	// Create options for projection
	opts := options.FindOne().SetProjection(projection)

	// Find the document by ID with projection and decode it into the provided variable
	err := collection.FindOne(ctx, filter, opts).Decode(doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrDocumentNotFound
		}
		return err
	}

	return nil
}

// Find retrieves documents from the collection based on the provided filter and projection.
// It takes a filter, projection, and a pointer to a slice where the results will be stored.
// It returns an error if the find operation fails.
func FindWithProjection(collection interface{}, ctx *gin.Context, filter bson.M, projection bson.M, results interface{}) error {

	// Create options for projection
	opts := options.Find().SetProjection(projection)

	// Find documents with projection and decode them into the provided results variable
	cur, err := collection.(*mongo.Collection).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	// Decode the results into the provided variable
	if err := cur.All(ctx, results); err != nil {
		return err
	}

	return nil
}

// Updates a Many documents in given collection
func UpdateMany(collection *mongo.Collection, ctx *gin.Context, filter interface{}, update interface{}) error {

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}

// BulkInsert inserts multiple documents into the collection using InsertMany.
// It takes a gin.Context, a slice of documents, and the collection where they should be inserted.
// It returns an error if the insert operation fails.
func BulkInsert(ctx *gin.Context, collection *mongo.Collection, documents []interface{}) error {
	// Perform the bulk insert operation
	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}
