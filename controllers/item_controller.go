package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tatrasoft/fyp-rest-backend-service/db"
	"github.com/tatrasoft/fyp-rest-backend-service/model"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

const (
	itemsCollection = "items"
)

func CreateItem(response http.ResponseWriter, request *http.Request) {
	item := &model.Item{}
	item.CreatedAt = time.Now().Unix()
	item.UpdatedAt =  time.Now().Unix()

	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}
	collection, err := db.GetDBCollection(dbname, itemsCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to insert the entry: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(result)
}

func GetItem(response http.ResponseWriter, request *http.Request) {
	request.Header.Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var item model.Item
	collection, err := db.GetDBCollection(dbname, itemsCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, model.Item{ID: id,}).Decode(&item)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to find the entry: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(item)
}

func GetItems(response http.ResponseWriter, request *http.Request) {
	var items []model.Item
	collection, err := db.GetDBCollection(dbname, itemsCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to retrieve the cursor: %v", err),
		})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		item :=  model.Item{}
		cursor.Decode(&item)
		items = append(items, item)
	}
	if err := cursor.Err(); err != nil {
		//
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(items)
}

func DeleteItem(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection, err := db.GetDBCollection(dbname, itemsCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to delete the record: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func UpdateItem(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var item model.Item

	// get id
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection, err := db.GetDBCollection(dbname, itemsCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	// read update model
	err = json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}

	// prepare model update
	update := bson.D{
		{
			"$set", bson.D{
			{"item_name", item.ItemName},
			{"item_description", item.ItemDescription},
			{"item_type", item.ItemType},
			{"item_img", item.ItemImg},
			{"item_price", item.ItemPrice},
			{"estimate_prepare_time", item.EstimatePrepareTime},
			{"updated_at", time.Now().Unix()},

		},
		},
	}

	sr := collection.FindOneAndUpdate(ctx, filter, update)
	if sr.Err() != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to update the record: %v", err),
		})
		return
	}
	item.ID = id

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(item)
}
