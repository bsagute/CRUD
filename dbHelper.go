package dbhelper

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetMongoClient To get mongo db client
func GetMongoClient(host, port, dbName, userName, password string, passwordSet bool) (*mongo.Client, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var client *mongo.Client
	var err error
	if passwordSet {
		client, err = mongo.Connect(ctx,
			options.Client().SetAuth(options.Credential{
				Username:    userName,
				Password:    password,
				AuthSource:  dbName,
				PasswordSet: passwordSet,
			}).ApplyURI("mongodb://@"+host+":"+port))
	} else {
		client, err = mongo.Connect(ctx,
			options.Client().ApplyURI("mongodb://@"+host+":"+port))
	}
	return client, ctx, err
}

// func GetMongoClient(host, port, dbName, userName, password string, passwordSet bool) (*mongo.Client, context.Context, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
// 	client, err := mongo.Connect(ctx,
// 		options.Client().SetAuth(options.Credential{
// 			Username:    userName,
// 			Password:    password,
// 			AuthSource:  dbName,
// 			PasswordSet: passwordSet,
// 		}).ApplyURI("mongodb://@"+host+":"+port))
// 	return client, ctx, err
// }

//GetMongoDB To Get the mongo db instance
func GetMongoDB(host, port, dbName, userName, password string, passwordSet bool) (*mongo.Database, context.Context, error) {
	client, ctx, err := GetMongoClient(host, port, dbName, userName, password, passwordSet)
	if err != nil {
		log.Print("Error While Connecting Mongo Client")
		return nil, ctx, err
	}
	return client.Database(dbName), ctx, nil
}
