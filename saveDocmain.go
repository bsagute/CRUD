package logic

import (
	"fmt"
	"net/http"
	"time"

	"digi-data-ingestion-client/clients/mongo"
	"digi-data-ingestion-client/models"
	"digi-data-ingestion-client/utils/logger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UpdateDocData = mongo.UpdateOne
var FindDocData = mongo.FindOne

func SaveCustomerDocuments(ctx *gin.Context, data []models.DocumentData, userRef string, CustomerId, DeviceId primitive.ObjectID, AndroidId string) (response models.SaveCustomerDocResponse) {
	newDocCount := 0

	query := bson.M{
		"_id": CustomerId,
		"phones": bson.M{
			"$elemMatch": bson.M{
				"android_id": AndroidId,
				"device_id":  DeviceId,
			},
		},
	}

	cursor := FindDocData(mongo.CustomersCollection, ctx, query, nil)

	if cursor.Err() != nil {
		response = models.SaveCustomerDocResponse{
			Message:      "error while finding data",
			ResponseCode: http.StatusInternalServerError,
			Success:      false,
			Data:         models.Data{},
		}
		return response
	}

	var result models.CustomersQuery
	if err := cursor.Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(err, "No documents found for the query.")
		}

		response = models.SaveCustomerDocResponse{
			Message:      "failed while decoding document details from DB",
			ResponseCode: http.StatusInternalServerError,
			Success:      false,
			Data:         models.Data{},
		}
		return response
	}

	var duplicateDocFound bool
	mapOfDoc := make(map[string]models.DocumentData)

	if len(result.Docs) != 0 {
		for _, v := range result.Docs {
			if v.Number != "" {
				mapOfDoc[v.Number] = v
			}
		}
	}

	var arrOfDoc []models.DocumentData
	for _, v := range data {
		if _, isFound := mapOfDoc[v.Number]; !isFound {
			mapOfDoc[v.Number] = v
			arrOfDoc = append(arrOfDoc, v)
			newDocCount++
		} else {
			logger.Info(fmt.Sprintf("Duplicate entry found for document with number: %s", v.Number))
			duplicateDocFound = true
			continue
		}
	}

	now := time.Now()
	date := now.Format("2006-01-02 15:04:05")

	if len(arrOfDoc) >= 1 {
		// Check for an empty document in the arrOfDoc slice
		isEmptyDoc := true
		for _, doc := range arrOfDoc {
			if !isEmptyDocument(doc) {
				isEmptyDoc = false
				break
			}
		}

		if isEmptyDoc {
			response = models.SaveCustomerDocResponse{
				Message:      "No docs found to insert",
				ResponseCode: 0,
				Success:      false,
				Data:         nil,
			}
			return response
		}

		filter := bson.M{"user_reference_number": userRef}
		update := bson.M{
			"$push": bson.M{"docs": bson.M{"$each": arrOfDoc}},
			"$set":  bson.M{"docs_last_updated_datetime": date},
		}

		err := UpdateDocData(mongo.CustomersCollection, ctx, filter, update)
		deviceUpdateFilter := bson.M{"_id": DeviceId}
		deviceUpdateUpdate := bson.M{"$set": bson.M{"docs_last_updated_datetime": date}}

		deviceUpdateErr := UpdateDocData(mongo.CustomerDevicesCollection, ctx, deviceUpdateFilter, deviceUpdateUpdate)
		if deviceUpdateErr != nil {
			logger.Error(deviceUpdateErr, "Error updating document details in DB")

			response = models.SaveCustomerDocResponse{
				Message:      "failed while updating docs_last_updated_datetime details in customer device collection",
				ResponseCode: http.StatusInternalServerError,
				Success:      false,
				Data:         models.Data{},
			}
			return response
		}
		if err != nil {
			logger.Error(err, "Error updating document details in DB")

			response = models.SaveCustomerDocResponse{
				Message:      "failed while updating document details in DB",
				ResponseCode: http.StatusInternalServerError,
				Success:      false,
				Data:         models.Data{},
			}
			return response
		}

		response = models.SaveCustomerDocResponse{
			Message:      "Docs saved successfully",
			ResponseCode: http.StatusOK,
			Success:      true,
			Data: models.Data{
				TotalAdded:              newDocCount,
				DocsLastUpdatedDatetime: date,
			},
		}
	} else if duplicateDocFound {
		response = models.SaveCustomerDocResponse{
			Message:      "Duplicate documents found.",
			ResponseCode: http.StatusOK,
			Success:      false,
			Data:         nil,
		}
	} else {
		response = models.SaveCustomerDocResponse{
			Message:      "Docs are required.",
			ResponseCode: http.StatusOK,
			Success:      false,
			Data:         nil,
		}
		return response
	}

	return response
}

// Function to check if a document is empty
func isEmptyDocument(doc models.DocumentData) bool {
	return (doc.Number == "" && doc.DocType == "")
}
