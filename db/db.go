package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client

func GetDBCollection(dbName string, dbCollection string) (*mongo.Collection, error) {
	collection := mongoClient.Database(dbName).Collection(dbCollection)

	return collection, nil
}

func SetMongoClient(client *mongo.Client) {
	mongoClient = client
}

